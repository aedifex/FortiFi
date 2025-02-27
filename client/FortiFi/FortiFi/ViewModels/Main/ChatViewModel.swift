//
//  ChatViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/23/25.
//

import Foundation

@MainActor final class ChatViewModel: ObservableObject {
    @Published var messages: [ChatMessage]
    @Published var input: String = ""
    @Published var isLoading: Bool = false
    private var threatId: Int?
    
    init() {
        let greeting = "Hello! I am your personal chatbot assistance specializing in IoT and Home network security. How may I assist you today?"
        messages = [
            ChatMessage(id: "startConversation", text: greeting, sender: 1)
        ]
    }
    
    init(for threatId: Int) {
        self.threatId = threatId
        let greeting = "Hello! I'm your network security assistant. Suspicious activity has been detected on the network and I'm here to help."
        messages = [
            ChatMessage(id: "startConversation", text: greeting, sender: 1)
        ]
    }
    
    func threatSpecified() -> Bool {
        return threatId != nil
    }
    
    func handleNeedRecommendationsResponse() async {
        if input.contains("yes") || input.contains("Yes") {
            input = ""
            await getRecommendations()
        } else {
            pushGoodbye()
        }
    }
    
    func pushGoodbye() {
        messages.append(ChatMessage(id: String(messages.count), text: "Goodbye! Feel free to reach out again if you need any help in the future. Have a great day!", sender: 1))
    }
    
    func pushUserMessage() {
        if input == "" { return }
        messages.append(ChatMessage(id: String(messages.count), text: input.trimmingCharacters(in: .whitespaces), sender: 0))
    }
    
    func getThreatAssistance() async {
        do{
            isLoading = true
            let resp = try await NetworkManager.shared.getThreatAssistance(threatId: threatId!)
            messages.append(resp)
            messages.append(
                ChatMessage(id: "offerRecommendations-\(messages.count)", text: "Do you want recommendations on how to resolve this issue? (yes/no)", sender: 1)
            )
            isLoading = false
        } catch{
            isLoading = false
            print("error in getThreatAssistance: \(error)")
        }
    }

    func offerMoreAssistance() {
        messages.append(
            ChatMessage(id: "furtherAssistance-\(messages.count)", text: "Is there anything else I can help you with?", sender: 1)
        )
    }
    
    func getRecommendations() async {
        do{
            isLoading = true
            let resp = try await NetworkManager.shared.getRecommendations(threatId: threatId!)
            messages.append(resp)
            offerMoreAssistance()
            isLoading = false
        }catch {
            isLoading = false
            print("error in getRecommendations: \(error)")
        }
    }
    
    func getMoreAssistance() async {
        if input.contains("no") || input.contains("No") {
            pushGoodbye()
            return
        }
        let question = input
        input = ""
        do{
            isLoading = true
            let resp = try await NetworkManager.shared.getMoreAssistance(threatId: threatId!, query: question)
            messages.append(resp)
            offerMoreAssistance()
            isLoading = false
        }catch {
            isLoading = false
            print("error in getMoreAssistance: \(error)")
        }
    }
    
    func getGeneralAssistance() async {
        let question = input
        input = ""
        do{
            isLoading = true
            let resp = try await NetworkManager.shared.getGeneralAssistance(query: question)
            messages.append(resp)
            offerMoreAssistance()
            isLoading = false
        }catch {
            isLoading = false
            print("error in getGeneralAssistance: \(error)")
        }
    }
    
}
