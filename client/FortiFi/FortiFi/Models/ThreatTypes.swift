//
//  ThreatTypes.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/21/25.
//

import Foundation

enum ThreatTypes: String, Codable {
    case benign = "0", portScan = "1", ddos = "2"
}
