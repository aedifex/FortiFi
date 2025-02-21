//
//  Event.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import Foundation

struct Event: Codable, Identifiable, Hashable{
    var threat_id: Int
    var id: String
    var details: String
    var ts: String
    var expires: String
    var type: ThreatTypes
    var src: String
    var dst: String
}
