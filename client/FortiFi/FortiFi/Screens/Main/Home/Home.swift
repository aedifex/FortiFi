//
//  Home.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/17/25.
//

import SwiftUI

struct Home: View {
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
                                .font(.headline)
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                        NavigationLink {
                            Text("Hello events")
                        } label: {
                            Label("Needs Attention", image: "warn")
                        }
                    }
                    
                    VStack (spacing: 16){
                        HStack {
                            Text("Insights")
                                .font(.headline)
                                .fontWeight(.regular)
                                .frame(maxWidth: .infinity, alignment: .leading)
                        }
                    }
                    
                }
                .padding()
                .background(Color("background"))
                
            }
        }
    }
}

#Preview {
    Home()
}
