//
//  HomeViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import Foundation

@MainActor final class HomeViewModel: ObservableObject {
    static let shared = HomeViewModel()
    
    @Published var events: [Event] = []
    
    func updateEvents() async {
        do {
            events = try await NetworkManager.shared.getEvents()
        } catch {
            print("error getting events: \(error)")
        }
    }
    
    func refresh() async {
        await updateEvents()
    }
    
}
