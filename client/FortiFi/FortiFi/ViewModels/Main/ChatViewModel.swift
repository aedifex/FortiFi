//
//  ChatViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/23/25.
//

import Foundation

@MainActor final class ChatViewModel: ObservableObject {
    @Published var messages: [ChatMessage] = []
    @Published var input: String = ""
    
    init() {
        messages = [
            .init(
                id: "123",
                text: "Good afternoon, Jonathan! What can I help you with?",
                sender: 1
            ),
            .init(
                id: "456",
                text: "What does this threat mean?",
                sender: 0
            ),
            .init(
                id: "789",
                text: "Certainly!\n\nThis is an example of a DDoS attack. It involves overwhelming a network or system with a large number of simultaneous requests or connections. Is there anything else I can assist you with?",
                sender: 1
            ),
        ]
    }
}
