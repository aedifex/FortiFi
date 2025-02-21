//
//  NetworkStatusNavigation.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct NetworkStatusNavigation: View {
    
    @ObservedObject var homeViewModel = HomeViewModel.shared

    var body: some View {
        
        if homeViewModel.threats.count == 0 {
            
            HStack {
                
                Image("OK")
                    .resizable()
                    .frame(width: 50, height: 50)
                
                VStack (alignment: .leading){
                    Text("Good")
                        .font(.body)
                        .foregroundColor(.fortifiForeground)
                    Text("^[\(homeViewModel.threats.count) Threat](inflect: true) found")
                        .Label()
                        .foregroundColor(.foregroundMuted)
                }
                
                Spacer()
                
            }
            .padding()
            .background(.fortifiBackground)
            .cornerRadius(12)
            .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
            
        } else {
            
            NavigationLink(destination: Events()) {
                
                HStack {
                    
                    Image("Error")
                        .resizable()
                        .frame(width: 50, height: 50)
                    VStack (alignment: .leading){
                        Text("Needs Attention")
                            .font(.body)
                            .foregroundColor(.fortifiForeground)
                        Text("^[\(homeViewModel.threats.count) Threat](inflect: true) found")
                            .Label()
                            .foregroundColor(.foregroundMuted)
                    }
                    
                    Spacer()
                    
                    Image(systemName: "chevron.right")
                        .foregroundColor(.foregroundMuted)
                }
                .padding()
                .background(.fortifiBackground)
                .cornerRadius(12)
                .shadow(color: Color.black.opacity(0.1), radius: 5, x: 0, y: 2)
                
            }
            .buttonStyle(PlainButtonStyle())
            
        }
    }
}

#Preview {
    NetworkStatusNavigation()
}
