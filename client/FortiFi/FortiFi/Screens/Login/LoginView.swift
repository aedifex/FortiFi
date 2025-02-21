//
//  LoginView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import SwiftUI

struct LoginView: View {
    @ObservedObject var viewModel = LoginViewModel.shared
    
    var body: some View {
        VStack(spacing: 28) {
            // Header
            VStack(spacing: 18) {
                Image("FortiFi-Logo")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
                    .frame(width: 120)
                Text("Welcome back to FortiFi. If you have not yet registered an account, login via the QR code on your FortiFi device.")
                    .font(.system(size: 14))
                    .foregroundColor(.gray)
                    .multilineTextAlignment(.center)
            }
            
            // Form Fields
            VStack(spacing: 24) {
                VStack(alignment: .leading, spacing: 8) {
                    Text("Email")
                        .font(.system(size: 14, weight: .medium))
                    TextField("", text: $viewModel.user.email)
                        .textFieldStyle(CustomTextFieldStyle())
                        .autocapitalization(.none)
                        .keyboardType(.emailAddress)
                }
                
                VStack(alignment: .leading, spacing: 8) {
                    Text("Password")
                        .font(.system(size: 14, weight: .medium))
                    SecureField("", text: $viewModel.user.password)
                        .textFieldStyle(CustomTextFieldStyle())
                }
            }
            
            // Login Button
            Button {
                Task {
                    await viewModel.loginUser()
                }
            } label: {
                    Text("Login")
                    .frame(maxWidth: .infinity)
                    .padding()
                    .background(.fortifiPrimary)
                    .foregroundStyle(.fortifiBackground)
                    .cornerRadius(8)
            }
            
        }
        .padding(40)
        .frame(maxHeight: .infinity)
        .background(.backgroundAlt)
        .alert(item: $viewModel.alert) {alert in
            Alert(title: alert.title, message: alert.message, dismissButton: alert.dismissButton)
        }

    }
}

struct CustomTextFieldStyle: TextFieldStyle {
    func _body(configuration: TextField<Self._Label>) -> some View {
        configuration
            .padding(.horizontal, 16)
            .padding(.vertical, 12)
            .background(.fortifiBackground)
            .foregroundColor(.gray)
            .cornerRadius(8)
            .overlay(
                RoundedRectangle(cornerRadius: 8)
                    .stroke(.foregroundMuted, lineWidth: 1)
            )
    }
}

#Preview {
    LoginView()
}
