//
//  DevicesViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/21/25.
//

import Foundation

@MainActor final class DevicesViewModel: ObservableObject {
    @Published var devices: [DevicesResponse] = []
    
    func getDevices() async {
        do {
            devices = try await NetworkManager.shared.getDevices()
        } catch {
            print("Error getting devices: \(error)")
        }
    }
    
}
