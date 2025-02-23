//
//  ContentView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/6/25.
//

import SwiftUI

struct ContentView: View {
    
    @ObservedObject var manager = BaseViewModel.shared
    
    var body: some View {
        switch manager.authenticated {
        case true:
            TabView{
                Home()
                    .tabItem{
                        SwiftUI.Label("Summary", systemImage: "house")
                    }
                    .toolbarBackground(.fortifiBackground, for: .tabBar)
                    .toolbarBackgroundVisibility(.visible, for: .tabBar)

                Devices()
                    .tabItem {
                        SwiftUI.Label("Devices", systemImage: "wifi.router")
                    }
                    .toolbarBackground(.fortifiBackground, for: .tabBar)
                    .toolbarBackgroundVisibility(.visible, for: .tabBar)

                Chat()
                    .tabItem {
                        SwiftUI.Label("Chat", systemImage: "bubble.left.and.exclamationmark.bubble.right")
                    }
                    .toolbarBackground(.fortifiBackground, for: .tabBar)
                    .toolbarBackgroundVisibility(.visible, for: .tabBar)
            }
            .tint(.fortifiPrimary)
            
        case false:
            LoginView()
        }
    }
}


#Preview {
    ContentView()
}
