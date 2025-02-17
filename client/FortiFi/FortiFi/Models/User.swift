//
//  User.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation

struct User: Codable {
    var email: String = ""
    var password: String = ""
    var first_name: String = ""
    var last_name: String = ""
    var id: String = ""
}
