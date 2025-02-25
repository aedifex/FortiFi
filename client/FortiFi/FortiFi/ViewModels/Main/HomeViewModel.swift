//
//  HomeViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import Foundation

@MainActor final class HomeViewModel: ObservableObject {
    static let shared = HomeViewModel()
    
    @Published var threats: [Event] = []
    @Published var eventCounts: DistributionResponse = DistributionResponse(Benign: 0, PortScan: 0, DDoS: 0, PrevWeekTotal: 0)
    @Published var totalEvents: Int = 0
    @Published var difference = 0
    @Published var distributions: [ThreatTypes: Double] = [
        .benign: 0.0,
        .portScan: 0.0,
        .ddos: 0.0
    ]
    
    func updateEvents() async {
        do {
            threats = try await NetworkManager.shared.getEvents()
        } catch {
            print("error getting events: \(error)")
        }
    }
    
    func getEventsDistribution() async {
        do {
            eventCounts = try await NetworkManager.shared.getEventsDistribution()
            totalEvents = eventCounts.Benign + eventCounts.PortScan + eventCounts.DDoS
            difference = totalEvents - eventCounts.PrevWeekTotal
            if totalEvents > 0 {
                distributions[.benign] = (Double(eventCounts.Benign) / Double(totalEvents)) * 100
                distributions[.portScan] = (Double(eventCounts.PortScan) / Double(totalEvents)) * 100
                distributions[.ddos] = (Double(eventCounts.DDoS) / Double(totalEvents)) * 100
            } else {
                distributions = [.benign: 0.0, .portScan: 0.0, .ddos: 0.0]
            }
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
