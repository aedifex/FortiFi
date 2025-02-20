package database

import (
	"database/sql"
	"fmt"

	"github.com/aedifex/FortiFi/config"
	"github.com/go-sql-driver/mysql"
)

type DatabaseConn struct {
    conn *sql.DB
}

const (
    UsersTable       = "Users"
    EventsTable      = "NetworkEvents"
    UserRefreshTable =  "UserRefreshTokens"
    PiRefreshTable   =  "PiRefreshTokens"
)

func ConnectDatabase(config *config.Config) (*DatabaseConn, error) {
	sqlConfig := mysql.Config{
        User:   config.DB_USER,
        Passwd: config.DB_PASS,
        Net:    "tcp",
        Addr:   config.DB_URL,
        DBName: config.DB_NAME,
    }

    // Get a database handle.
    db, err := sql.Open("mysql", sqlConfig.FormatDSN())
    if err != nil {
        return nil, fmt.Errorf("error opening the database: %s", err.Error())
    }

    err = db.Ping()
    if err != nil {
        return nil, fmt.Errorf("could not connect to DB: %s", err.Error())
    }
    
	return &DatabaseConn{
        conn: db,
    }, nil

}

func (db *DatabaseConn) InsertUser(user *User) (*DatabaseError) {
    //Check if user exists
    userExists := db.userExists(user)
    if userExists != nil {
        return USER_EXISTS_ERROR
    }
    
    query := fmt.Sprintf("INSERT into %s (id,first_name,last_name,email,password) VALUES (?,?,?,?,?);", UsersTable)

    // Hash the password
    hashedPassword, err := HashString(user.Password)
    if err != nil {
        return HASH_ERROR(err)
    }

    // prepare statement
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    // execute statement
    _, err = preparedStatement.Exec(user.Id, user.FirstName, user.LastName, user.Email, hashedPassword)
    if err != nil {
        return EXEC_ERROR(err)
    }

    return nil
}

func (db *DatabaseConn) UpdateFcmToken(subjectId string, fcmToken string) *DatabaseError {
    
    userIdExists, existsErr := db.userIdExists(subjectId)
    if existsErr != nil {
        return existsErr
    }
    if !userIdExists {
        return DNE_ERROR
    }

    // insert new fcm token
    query := fmt.Sprintf("UPDATE %s SET fcm_token = ? WHERE id = ?;", UsersTable)

    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }

    // execute statement
    _, err = preparedStatement.Exec(fcmToken, subjectId)
    if err != nil {
        return QUERY_ERROR(err)
    }

    return nil
}

// Takes in a user object and compares the passwords to the matching user in the database.
//
// Returns a found user, http status code, and nil if the user is validated.
// Otherwise returns nil, error code, and a corresponding error.
func (db *DatabaseConn) ValidateLogin(user *User) (*User, *DatabaseError) {

    query := fmt.Sprintf("SELECT id, first_name, last_name, email, password FROM %s WHERE email = ?;", UsersTable)

    // create prepared statment
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return nil, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    // input validation
    if user.Email == "" || user.Password == "" {
        return nil, INVALID_INPUT_ERROR
    }

    // Query database
    res, err := preparedStatement.Query(user.Email)
    if err != nil {
        return nil, QUERY_ERROR(err)
    }

    // check if user exists
    if !res.Next() { return nil, DNE_ERROR }

    // store the user the query returned
    foundUser := &User{}
    if err := res.Scan(&foundUser.Id, &foundUser.FirstName, &foundUser.LastName, &foundUser.Email, &foundUser.Password); err != nil {
        return nil, SCAN_ERROR(err)
    }

    // validate the user
    if !HashMatch(foundUser.Password, user.Password) {
        return nil, UNAUTHORIZED_ERROR
    }
    return foundUser, nil
}

// Validates a refresh token by finding a match in the database.
// Returns error if token is invalid, otherwise nil.
func (db *DatabaseConn) ValidateRefresh(refresh *RefreshToken, tablename string) *DatabaseError {
    
    // drop expired tokens from respective table
    query := fmt.Sprintf("DELETE FROM %s WHERE expires < Now();", tablename)
    _, err := db.conn.Exec(query)
    if err != nil {
        return EXEC_ERROR(err)
    }
    
    // check if matching refresh token is in the database
    query = fmt.Sprintf("SELECT token_hash FROM %s WHERE id = ?;",tablename)

    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }

    rows, err := preparedStatement.Query(refresh.Id)
    if err != nil {
        return QUERY_ERROR(err)
    }

    if !rows.Next() {
        return DNE_ERROR
    }

    // Get token info
    storedHashedToken := ""
    err = rows.Scan(&storedHashedToken)
    if err != nil {
        return SCAN_ERROR(err)
    }
    
    if !HashMatch(storedHashedToken, refresh.Token) {
        return UNAUTHORIZED_ERROR
    }

    return nil
}

// Stores the hash of a given token and subject id with exp time 7 days from current time the given tablename.
// Returns nil on success or non-nil error on failure.
func (db *DatabaseConn) StoreRefresh(refresh *RefreshToken, tablename string) *DatabaseError {

    // Delete existing keys for this user
    query := fmt.Sprintf("DELETE FROM %s WHERE id = ?;", tablename)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(refresh.Id)
    if err != nil {
        return EXEC_ERROR(err)
    }
    
    // Insert new token into database
    query = fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?);", tablename) // token, id, expires
    preparedStatement, err = db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    tokenHash, err := HashString(refresh.Token)
    if err != nil {
        return HASH_ERROR(err)
    }
    _, err = preparedStatement.Exec(tokenHash, refresh.Id, refresh.Expires)
    if err != nil {
        return EXEC_ERROR(err)
    }
    return nil
}


// Checks if a given user exists in the database.
// Returns non-nil if the user email or id exists, otherwise.
func (db *DatabaseConn) userExists(user *User) *DatabaseError {
    // check email
    queryEmail := fmt.Sprintf("SELECT * FROM %s WHERE email = ?;", UsersTable)
    preparedStatementEmail, _ := db.conn.Prepare(queryEmail)
    defer preparedStatementEmail.Close()
    res, err := preparedStatementEmail.Query(user.Email)
    if err != nil { return EXEC_ERROR(err) }
    defer res.Close()
    if res.Next() { return USER_EXISTS_ERROR}
    
    // check id
    queryId := "SELECT * FROM USERS WHERE id = ?;"
    preparedStatementId, _ := db.conn.Prepare(queryId)
    defer preparedStatementId.Close()
    res, err = preparedStatementId.Query(user.Id)
    if err != nil { return QUERY_ERROR(err) }
    defer res.Close()
    if res.Next() { return USER_EXISTS_ERROR }

    return nil
}

// Checks if a given user id exists in the database.
// Returns true if the user id exists, otherwise false.
func (db *DatabaseConn) userIdExists(id string) (bool, *DatabaseError) {

    query := fmt.Sprintf("SELECT * FROM %s WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return false, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    res, err := preparedStatement.Query(id)
    if err != nil { return false, QUERY_ERROR(err) }
    defer res.Close()
    if res.Next() { return true, nil }

    return false, nil

}

func (db *DatabaseConn) StoreEvent(e *Event) *DatabaseError {

    // Insert id, details, ts, expires
    query := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?)", EventsTable)

    // prepare statement
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }

    // execute statement
    res, err := preparedStatement.Exec(e.Id, e.Details, e.TS, e.Expires, e.Type)
    if err != nil {
        return EXEC_ERROR(err)
    }

    // check rows affected
    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return ROWS_AFFECTED_ERROR(err)
    }
    if rowsAffected == 0 {
        return DNE_ERROR
    }
    
    return nil
}

func (db *DatabaseConn) GetFcmToken(subjectId string) (string, *DatabaseError) {
    
    // prepare query
    query := fmt.Sprintf("SELECT fcm_token from %s where id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return "", PREPARE_ERROR(err)
    }

    // execute query
    rows, err := preparedStatement.Query(subjectId)
    if err != nil {
        return "", QUERY_ERROR(err)
    }

    if !rows.Next() {
        return "", DNE_ERROR
    }

    fcmToken := ""
    err = rows.Scan(&fcmToken)
    if err != nil {
        return "", SCAN_ERROR(err)
    }

    return fcmToken,nil
}

func (db *DatabaseConn) Close() error {
    return db.conn.Close()
}

func (db *DatabaseConn) GetUserEvents(userId string) ([]*Event, *DatabaseError) {

    // delete expired events
    deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE expires < NOW();", EventsTable)
    _, err := db.conn.Exec(deleteQuery)
    if err != nil {
        return nil, EXEC_ERROR(err)
    }

    // prepare query
    query := fmt.Sprintf("SELECT id, details, ts, expires, event_type FROM %s WHERE id = ? ORDER BY ts DESC;", EventsTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return nil, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    // execute query
    rows, err := preparedStatement.Query(userId)
    if err != nil {
        return nil, QUERY_ERROR(err)
    }
    defer rows.Close()

    var events []*Event
    for rows.Next() {
        event := &Event{}
        err := rows.Scan(&event.Id, &event.Details, &event.TS, &event.Expires, &event.Type)
        if err != nil {
            return nil, SCAN_ERROR(err)
        }
        events = append(events, event)
    }

    if err = rows.Err(); err != nil {
        return nil, QUERY_ERROR(err)
    }

    return events, nil
}

func (db *DatabaseConn) UpdateWeeklyDistribution(userId string, normal int, anomalous int, malicious int) *DatabaseError {

    userExists, userExistsErr := db.userIdExists(userId);
    if userExistsErr != nil {
        return userExistsErr
    }
    if !userExists {
        return DNE_ERROR
    }

    query := fmt.Sprintf("UPDATE %s SET normal_count = ?, anomalous_count = ?, malicious_count = ? WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(normal, anomalous, malicious, userId)
    if err != nil {
        return EXEC_ERROR(err)
    }

    return nil
}

func (db *DatabaseConn) GetWeeklyDistribution(userId string) (*WeeklyDistribution, *DatabaseError) {

    userExists, userExistsErr := db.userIdExists(userId);
    if userExistsErr != nil {
        return nil, userExistsErr
    }
    if !userExists {
        return nil, DNE_ERROR
    }

    query := fmt.Sprintf("SELECT normal_count, anomalous_count, malicious_count FROM %s WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return nil, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    rows, err := preparedStatement.Query(userId)
    if err != nil { return nil, QUERY_ERROR(err) }
    defer rows.Close()

    if !rows.Next() {
        return nil, DNE_ERROR
    }   

    weeklyDistribution := &WeeklyDistribution{}
    err = rows.Scan(&weeklyDistribution.Normal, &weeklyDistribution.Anomalous, &weeklyDistribution.Malicious)
    if err != nil {
        return nil, SCAN_ERROR(err)
    }
    
    return weeklyDistribution, nil
}

