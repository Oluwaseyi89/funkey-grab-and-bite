// Builds a CSV file and starts the download. Values that start with a
// formula character are prefixed so spreadsheets don't run them.
export const downloadCsv = (filename: string, header: string[], rows: (string | number | null | undefined)[][]) => {
  const escape = (value: string | number | null | undefined) => {
    let text = value == null ? '' : String(value);
    if (/^[=+\-@]/.test(text)) text = `'${text}`;
    return /[",\n\r]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
  };
  const csv = [header, ...rows].map((row) => row.map(escape).join(',')).join('\r\n');
  const blob = new Blob([`﻿${csv}`], { type: 'text/csv;charset=utf-8' });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  setTimeout(() => URL.revokeObjectURL(url), 1000);
};
