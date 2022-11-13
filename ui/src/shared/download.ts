function formatDownloadName(filename: string, extension: string) {
  return `${filename.trim().toLowerCase().replace(/\s/g, '_')}${extension}`
}

export const downloadTextFile = (
  text: string,
  filename: string,
  extenstion: string,
  mimeType: string = 'text/plain'
) => {
  const formattedName = formatDownloadName(filename, extenstion)
  const blob = new Blob([text], {type: mimeType})
  const a = document.createElement('a')

  a.href = window.URL.createObjectURL(blob)
  a.target = '_blank'
  a.download = formattedName

  document.body.appendChild(a)
  a.click()
  // eslint-disable-next-line @typescript-eslint/no-non-null-assertion
  a.parentNode!.removeChild(a)
}
