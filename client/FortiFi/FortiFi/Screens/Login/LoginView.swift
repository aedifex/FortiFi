//
//  LoginView.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/16/25.
//

import SwiftUI

struct LoginView: View {
    @ObservedObject var viewModel = LoginViewModel()
    var body: some View {
        VStack(spacing: 32) {
            // Header
            VStack(spacing: 15) {
                Image("FortiFi-Logo")
                    .resizable()
                    .aspectRatio(contentMode: .fit)
                    .frame(width: 150)
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
                    .background(Color("Fortifi-Primary"))
                    .foregroundColor(.white)
                    .cornerRadius(8)
            }
            
        }
        .padding(40)
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
            .background(Color.white)
            .foregroundColor(.gray)
            .cornerRadius(8)
            .overlay(
                RoundedRectangle(cornerRadius: 8)
                    .stroke(Color("Foreground-Muted"), lineWidth: 1)
            )
    }
}

#Preview {
    LoginView()
}
