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
    @Published var distribution: DistributionResponse = DistributionResponse(Normal: 0, Anomalous: 0, Malicious: 0, PrevCount: 0)
    @Published var totalEvents: Int = 0
    @Published var difference = 0
    
    func updateEvents() async {
        do {
            events = try await NetworkManager.shared.getEvents()
        } catch {
            print("error getting events: \(error)")
            BaseViewModel.shared.authenticated = false
        }
    }
    
    func getEventsDistribution() async {
        do {
            distribution = try await NetworkManager.shared.getEventsDistribution()
            totalEvents = distribution.Anomalous + distribution.Normal + distribution.Malicious
            difference = totalEvents - distribution.PrevCount
        }
        catch {
            print("error getting distribution info: \(error)")
            BaseViewModel.shared.authenticated = false
        }
    }
    
    func refresh() async {
        await updateEvents()
        await getEventsDistribution()
    }
    
}
