//
//  LoginViewModel.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import Foundation

@MainActor final class LoginViewModel: ObservableObject {
    
    static var shared = LoginViewModel()
    
    @Published var user = User()
    @Published var isLoading = false
    @Published var alert: AlertItem?
    
    func loginUser() async{
        
        isLoading = true

        do {
            try await NetworkManager.shared.login(user)
        } catch {
            switch error {
            case Errors.inputError:
                alert = AlertContext.inputError
            case Errors.internalError:
                alert = AlertContext.general
            case Errors.networkError:
                alert = AlertContext.networkError
            case Errors.invalidUrl:
                alert = AlertContext.general
            case Errors.notFound:
                alert = AlertContext.unauthorized
            case Errors.unauthorized:
                alert = AlertContext.unauthorized
            default:
                alert = AlertContext.general
            }
        }
        isLoading = false
        
    }
    
    
}
