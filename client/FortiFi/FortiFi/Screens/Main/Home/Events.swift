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
                .Title()
                .fontWeight(.medium)
                .frame(maxWidth: .infinity, alignment: .leading)
            Text("\(viewModel.threats.count) Threats found")
                .Label()
                .foregroundColor(.foregroundMuted)
            
            ScrollView{
                VStack (spacing: 12) {
                    ForEach(viewModel.threats, id: \.threat_id){ threat in
                        EventTab(threat: threat)
                        
                        if threat.self != viewModel.threats.last.self {
                                   Divider()
                        }
                    }
                }
                .padding()
                .background(.fortifiBackground)
                .cornerRadius(16)
                .shadow(color: Color.black.opacity(0.1), radius: 5, x: 2, y: 2)
            }
            .contentMargins(2)
        }
        .padding()
        .toolbarBackground(.fortifiBackground, for: .navigationBar)
        .frame(maxHeight: .infinity)
        .background(.backgroundAlt)
        .foregroundStyle(.fortifiForeground)
    }
}

#Preview {
    Events()
}
