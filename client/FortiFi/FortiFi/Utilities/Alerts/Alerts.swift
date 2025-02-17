//
//  Alerts.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import SwiftUI

struct AlertItem: Identifiable {
    let id = UUID()
    let title: Text
    let message: Text
    let dismissButton: Alert.Button
}


struct AlertContext {
    
    static let networkError = AlertItem(title: Text("Network Error"),
                                      message: Text("Server Could not be Reached"),
                                      dismissButton: .default(Text("OK")))
    
    static let inputError = AlertItem(title: Text("Invalid Input"),
                                      message: Text("Please check your inputs"),
                                      dismissButton: .default(Text("OK")))
    
    static let notFound = AlertItem(title: Text("No User Found"),
                                      message: Text("No accounts match our records"),
                                      dismissButton: .default(Text("OK")))
    
    static let unauthorized = AlertItem(title: Text("Invalid Login"),
                                      message: Text("Login information is incorrect"),
                                      dismissButton: .default(Text("OK")))
    
    static let expiredToken = AlertItem(title: Text("Expired Session"),
                                        message: Text("Your session has expired. Please log back in."),
                                        dismissButton: .default(Text("OK")))
    
    static let general = AlertItem(title: Text("Error"),
                                  message: Text("An error has occured. Please try again."),
                                  dismissButton: .default(Text("OK")))
}
