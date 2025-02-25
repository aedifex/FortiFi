//
//  ChatbotResponses.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/25/25.
//

import Foundation

struct ChatResponse: Codable {
    var id: String
    var response: String
}

struct ChatMessage: Codable, Identifiable {
    var id: String
    var text: String
    var sender: Int = 1
}

