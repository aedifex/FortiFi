//
//  Date+.swift
//  FortiFi
//
//  Created by Jonathan Nguyen on 2/20/25.
//

import Foundation
// Taken from https://www.reddit.com/r/swift/comments/f8ai10/finding_start_and_end_of_current_week_gives_me/
extension Calendar {
    /*
    Week boundary is considered the start of
    the first day of the week and the end of
    the last day of the week
    */
    typealias WeekBoundary = (startOfWeek: Date?, endOfWeek: Date?)
    
    func currentWeekBoundary() -> WeekBoundary? {
        return weekBoundary(for: Date())
    }
    
    func weekBoundary(for date: Date) -> WeekBoundary? {
        let components = dateComponents([.yearForWeekOfYear, .weekOfYear], from: date)
        
        guard let startOfWeek = self.date(from: components) else {
            return nil
        }
        
        let endOfWeekOffset = weekdaySymbols.count - 1
        let endOfWeekComponents = DateComponents(day: endOfWeekOffset, hour: 23, minute: 59, second: 59)
        guard let endOfWeek = self.date(byAdding: endOfWeekComponents, to: startOfWeek) else {
            return nil
        }
        
        return (startOfWeek, endOfWeek)
    }
}
