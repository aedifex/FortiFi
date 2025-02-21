//
//  Events.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/19/25.
//

import SwiftUI

struct Events: View {
    @ObservedObject var viewModel = HomeViewModel.shared
    var body: some View {
        VStack(alignment: .leading, spacing: 24) {
            Text("Flagged Events")
                .font(.title)
                .fontWeight(.medium)
                .frame(maxWidth: .infinity, alignment: .leading)
            Text("\(viewModel.threats.count) Threats found")
                .font(.subheadline)
                .foregroundColor(Color("Foreground-Muted"))
            
            ScrollView{
                VStack (spacing: 12) {
                    ForEach(viewModel.threats, id: \.self){ threat in
                        EventTab(threat: threat)
                        
                        if threat.self != viewModel.threats.last.self {
                                   Divider()
                        }
                    }
                }
                .padding()
                .background(Color("Background"))
                .cornerRadius(16)
                .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
            }
            .contentMargins(2)
        }
        .padding()
        .toolbarBackground(Color("Background"), for: .navigationBar)
        .frame(maxHeight: .infinity)
        .background(Color("Background"))
        .foregroundStyle(Color("Foreground"))
    }
}

#Preview {
    Events()
}
