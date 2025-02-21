//
//  ContentView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/6/25.
//

import SwiftUI

struct ContentView: View {
    
    @ObservedObject var manager = BaseViewModel.shared
    @State var selection = "home"
    
    var body: some View {
        switch manager.authenticated {
        case true:
            TabView(selection: $selection){
                Home()
                    .tabItem {
                        if selection == "home" {
                            Image("home-active")
                        } else {
                            Image("home")
                        }
        
                    }
                    .tag("home")
                    .toolbarBackground(.white, for: .tabBar)
                    

                Devices()
                    .tabItem {
                        if selection == "devices" {
                            Image("devices-active")
                        } else {
                            Image("devices")
                        }
        
                    }
                    .tag("devices")

                Chat()
                    .tabItem {
                        if selection == "chat" {
                            Image("chatbot-active")
                        } else {
                            Image("chatbot")
                        }
        
                    }
                    .tag("chat")

            }
            .onAppear{
                Task {
                    await HomeViewModel.shared.refresh()
                }
            }
        case false:
            LoginView()
        }
    }
}


#Preview {
    ContentView()
}
