//
//  Base.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/17/25.
//

import Foundation

@MainActor final class BaseViewModel: ObservableObject {
    static var shared = BaseViewModel()
    @Published var loginSuccess = false
    
}
