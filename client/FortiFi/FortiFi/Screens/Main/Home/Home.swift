//
//  Home.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/17/25.
//

import SwiftUI

struct Home: View {
    
    @ObservedObject var homeViewModel = HomeViewModel.shared
    
    var body: some View {
        NavigationStack{
            ScrollView{
                VStack (alignment: .leading, spacing: 32){
                    
                    HStack {
                        Text("Summary")
                            .Title()
                            .fontWeight(.medium)
                            .frame(maxWidth: .infinity, alignment: .leading)
                    }
                    
                    VStack (spacing: 16){
                        HStack {
                            Text("Network Status")
                                .Header()
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        NetworkStatusNavigation()
                    }
                    
                    VStack (spacing: 16){
                        HStack {
                            Text("Event Insights")
                                .Header()
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        TotalEventsTab()
                        EventsDistribution()
                    }
                    
                }
            }
            .contentMargins(16)
            .frame(maxHeight: .infinity)
            .background(.backgroundAlt)
            .foregroundStyle(.fortifiForeground)
        }
        .onAppear{
            Task {
                await HomeViewModel.shared.refresh()
            }
        }
        .refreshable {
            Task {
                await homeViewModel.refresh()
            }
        }
    }
}

#Preview {
    Home()
}
