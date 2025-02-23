//
//  DevicesResponse.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/21/25.
//

import Foundation

struct DevicesResponse: Codable {
    var id: Int
    var name: String
    var ip_address: String
    var mac_address: String
}
