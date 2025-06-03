/**
 * Formats a date from database to YYYY-MM-DD HH:MM format
 * @param {string} dbTimeString - Time string from database (created_at field)
 * @returns {string} Formatted time string in YYYY-MM-DD HH:MM format
 */
export function formatTime(dbTimeString) {
    if (!dbTimeString) {
        return 'No date';
    }
    
    // Create date object from database time string
    const date = new Date(dbTimeString);
    
    // Check if date is valid
    if (isNaN(date.getTime())) {
        return 'Invalid Date';
    }
    
    // Format: YYYY-MM-DD HH:MM
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    
    return `${year}-${month}-${day} ${hours}:${minutes}`;
}
