//
//  LoginResponse.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation

struct LoginResponse: Codable {
    var jwt: String
    var refresh: String
}
