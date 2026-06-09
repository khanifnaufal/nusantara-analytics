import { isRef, type Ref } from 'vue'

/**
 * useExportCsv
 * Composable to convert an array of objects to a CSV file and trigger a browser download.
 * 
 * Supports both binding data/filename at initialization, or passing them to the exportCsv function.
 * 
 * @param initialData - Optional array or ref of objects to export
 * @param initialFilename - Optional filename or ref for the exported file
 */
export const useExportCsv = (
  initialData?: any[] | Ref<any[]>,
  initialFilename?: string | Ref<string>
) => {
  const exportCsv = (
    overrideData?: any[],
    overrideFilename?: string
  ) => {
    const dataVal = overrideData || (isRef(initialData) ? initialData.value : initialData)
    const filenameVal = overrideFilename || (isRef(initialFilename) ? initialFilename.value : initialFilename) || 'export.csv'

    if (!dataVal || !Array.isArray(dataVal) || dataVal.length === 0) {
      console.warn('useExportCsv: No data to export')
      return
    }

    // Convert array of objects to CSV string
    const headers = Object.keys(dataVal[0])
    const csvContent = [
      headers.join(','),
      ...dataVal.map(row => 
        headers.map(fieldName => {
          const value = row[fieldName]
          let stringVal = value === null || value === undefined ? '' : String(value)
          
          // Escape quotes and wrap in quotes if contains comma, newline, or quotes
          if (stringVal.includes(',') || stringVal.includes('"') || stringVal.includes('\n') || stringVal.includes('\r')) {
            stringVal = `"${stringVal.replace(/"/g, '""')}"`
          }
          return stringVal
        }).join(',')
      )
    ].join('\n')

    // Trigger download via blob URL
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', filenameVal.endsWith('.csv') ? filenameVal : `${filenameVal}.csv`)
    link.style.visibility = 'hidden'
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
  }

  return {
    exportCsv
  }
}
