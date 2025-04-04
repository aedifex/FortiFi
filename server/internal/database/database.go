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
    UsersTable       =  "Users"
    EventsTable      =  "NetworkThreats"
    UserRefreshTable =  "UserRefreshTokens"
    PiRefreshTable   =  "PiRefreshTokens"
    DevicesTable     =  "Devices"
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
    defer preparedStatement.Close()

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
    defer res.Close()
    
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
    defer preparedStatement.Close()

    rows, err := preparedStatement.Query(refresh.Id)
    if err != nil {
        return QUERY_ERROR(err)
    }
    defer rows.Close()

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
    preparedStatementEmail, err := db.conn.Prepare(queryEmail)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatementEmail.Close()
    res, err := preparedStatementEmail.Query(user.Email)
    if err != nil { return EXEC_ERROR(err) }
    defer res.Close()
    if res.Next() { return USER_EXISTS_ERROR}
    
    // check id
    queryId := fmt.Sprintf("SELECT * FROM %s WHERE id = ?;", UsersTable)
    preparedStatementId, err := db.conn.Prepare(queryId)
    if err != nil {
        return PREPARE_ERROR(err)
    }
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

    userExists, userExistsErr := db.userIdExists(e.Id)
    if userExistsErr != nil {
        return userExistsErr
    }
    if !userExists {
        return DNE_ERROR
    }

    
    // Insert id, details, ts, expires
    query := fmt.Sprintf("INSERT INTO %s (id, details, ts, expires, event_type, src_ip, dst_ip, confidence_interval) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", EventsTable)

    // prepare statement
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    // execute statement
    res, err := preparedStatement.Exec(e.Id, e.Details, e.TS, e.Expires, e.Type, e.SrcIP, e.DstIP, e.Confidence)
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
    defer preparedStatement.Close()

    // execute query
    rows, err := preparedStatement.Query(subjectId)
    if err != nil {
        return "", QUERY_ERROR(err)
    }
    defer rows.Close()
    
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
    query := fmt.Sprintf("SELECT threat_id, id, details, ts, expires, event_type, src_ip, dst_ip FROM %s WHERE id = ? ORDER BY ts DESC;", EventsTable)
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
        err := rows.Scan(&event.ThreatId, &event.Id, &event.Details, &event.TS, &event.Expires, &event.Type, &event.SrcIP, &event.DstIP)
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

func (db *DatabaseConn) UpdateWeeklyDistribution(userId string, benign int, portScan int, ddos int) *DatabaseError {

    userExists, userExistsErr := db.userIdExists(userId);
    if userExistsErr != nil {
        return userExistsErr
    }
    if !userExists {
        return DNE_ERROR
    }

    query := fmt.Sprintf("UPDATE %s SET benign_count = ?, port_scan_count = ?, ddos_count = ? WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(benign, portScan, ddos, userId)
    if err != nil {
        return EXEC_ERROR(err)
    }

    return nil
}

func (db *DatabaseConn) ResetWeeklyDistribution(userId string, weekTotal int) *DatabaseError {

    userExists, userExistsErr := db.userIdExists(userId);
    if userExistsErr != nil {
        return userExistsErr
    }
    if !userExists {
        return DNE_ERROR
    }

    query := fmt.Sprintf("UPDATE %s SET benign_count = 0, port_scan_count = 0, ddos_count = 0, prev_week_total = ? WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(weekTotal, userId)
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

    query := fmt.Sprintf("SELECT benign_count, port_scan_count, ddos_count, prev_week_total FROM %s WHERE id = ?;", UsersTable)
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
    err = rows.Scan(&weeklyDistribution.Benign, &weeklyDistribution.PortScan, &weeklyDistribution.DDoS, &weeklyDistribution.PrevWeekTotal)
    if err != nil {
        return nil, SCAN_ERROR(err)
    }
    
    return weeklyDistribution, nil
}

func (db *DatabaseConn) GetDevices(userId string) ([]*Device, *DatabaseError) {

    query := fmt.Sprintf("SELECT id, name, ip_address, mac_address, date_added, incident_count FROM %s WHERE userId = ?;", DevicesTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return nil, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    rows, err := preparedStatement.Query(userId)
    if err != nil {
        return nil, QUERY_ERROR(err)
    }
    defer rows.Close()

    var devices []*Device
    for rows.Next() {
        device := &Device{}
        err := rows.Scan(&device.Id, &device.Name, &device.IpAddress, &device.MacAddress, &device.DateAdded, &device.IncidentCount)
        if err != nil {
            return nil, SCAN_ERROR(err)
        }
        devices = append(devices, device)
    }

    return devices, nil
}       

func (db *DatabaseConn) AddDevice(device *Device) *DatabaseError {

    userExists, userExistsErr := db.userIdExists(device.UserId)
    if userExistsErr != nil {
        return userExistsErr
    }
    if !userExists {
        return DNE_ERROR
    }
    
    query := fmt.Sprintf("INSERT INTO %s (name, ip_address, mac_address, userId, date_added) VALUES (?, ?, ?, ?, ?);", DevicesTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    _, err = preparedStatement.Exec(device.Name, device.IpAddress, device.MacAddress, device.UserId, device.DateAdded)
    if err != nil {
        return EXEC_ERROR(err)
    }

    return nil 
}

func (db *DatabaseConn) GetThreatById(threatId int, userId string) (*Event, *DatabaseError) {

    query := fmt.Sprintf("SELECT details, ts, event_type, src_ip, dst_ip, confidence_interval FROM %s WHERE threat_id = ? AND id = ?;", EventsTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return nil, PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()

    rows, err := preparedStatement.Query(threatId, userId)
    if err != nil {
        return nil, QUERY_ERROR(err)
    }
    defer rows.Close()

    if !rows.Next() {
        return nil, DNE_ERROR
    }
    
    threat := &Event{}
    err = rows.Scan(&threat.Details, &threat.TS, &threat.Type, &threat.SrcIP, &threat.DstIP, &threat.Confidence)
    if err != nil {
        return nil, SCAN_ERROR(err)
    }

    return threat, nil
}

func (db *DatabaseConn) DeleteUser(id string) (*DatabaseError) {
    query := fmt.Sprintf("DELETE * from %s WHERE id = ?;", UsersTable)
    preparedStatement, err := db.conn.Prepare(query)
    if err != nil {
        return PREPARE_ERROR(err)
    }
    defer preparedStatement.Close()
    preparedStatement.Exec(id)
    return nil
}
