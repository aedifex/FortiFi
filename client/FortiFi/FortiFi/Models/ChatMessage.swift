//
//  ChatMessage.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/23/25.
//

import Foundation

struct ChatMessage: Codable, Identifiable {
    var id: String
    var text: String
    var sender: Int
}
