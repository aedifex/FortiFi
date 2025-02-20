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
    @Published var eventCounts: DistributionResponse = DistributionResponse(Normal: 0, Anomalous: 0, Malicious: 0, PrevWeekTotal: 0)
    @Published var totalEvents: Int = 0
    @Published var difference = 0
    @Published var distributions = [0.0, 0.0, 0.0]
    
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
            eventCounts = try await NetworkManager.shared.getEventsDistribution()
            totalEvents = eventCounts.Anomalous + eventCounts.Normal + eventCounts.Malicious
            difference = totalEvents - eventCounts.PrevWeekTotal
            distributions[0] = (Double(eventCounts.Normal) / Double(totalEvents)) * 100
            distributions[1] = (Double(eventCounts.Anomalous) / Double(totalEvents)) * 100
            distributions[2] = (Double(eventCounts.Malicious) / Double(totalEvents)) * 100
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
