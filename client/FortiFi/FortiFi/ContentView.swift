//
//  ContentView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/6/25.
//

import SwiftUI

struct ContentView: View {
    
    @ObservedObject var appModel = BaseViewModel.shared
    var body: some View {
        switch appModel.loginSuccess{
        case true:
            Home()
        case false:
            LoginView()
        }
    }
}


#Preview {
    ContentView()
}
