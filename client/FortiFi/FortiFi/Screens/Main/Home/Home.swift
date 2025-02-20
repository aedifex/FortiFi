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
                            .font(.title)
                            .fontWeight(.medium)
                            .frame(maxWidth: .infinity, alignment: .leading)
                    }
                    
                    VStack (spacing: 16){
                        HStack {
                            Text("Network Status")
                                .font(.title3)
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        NetworkStatusNavigation()
                    }
                    
                    VStack (spacing: 16){
                        HStack {
                            Text("Events Insights")
                                .font(.title3)
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        TotalEventsTab()
                    }
                    
                }
                .padding()
                .background(Color("Background"))
                
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
