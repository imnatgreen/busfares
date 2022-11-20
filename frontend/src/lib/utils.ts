export const apiUrl = (path: string) => `${import.meta.env.VITE_API_URL || ''}${path}`;
export const otpBase = 'https://otp.nat.omg.lol/otp';

export function timestampToString(timestamp: number) {
  const dt = new Date(timestamp);
  const time = dt.toLocaleTimeString('en-GB', {timeStyle: 'short'});
  const date = dt.toLocaleDateString('sv-SE', {dateStyle: 'short'});
  return {time, date};
};

// only works up to 24 hours
export function secondsToDurationString(seconds) {
  const duration = new Date(seconds*1000)
  if (duration.getUTCHours() + (duration.getUTCMinutes()/60) >= 1) {
    return `${
      duration.getUTCHours()
    } hr${duration.getUTCHours()>1?'s':''} ${duration.getUTCMinutes()} min${duration.getUTCMinutes()>1?'s':''}`;
  }
  if (duration.getUTCMinutes() < 1) {
    return `< 1 min`;
  }

  return `${duration.getUTCMinutes()} min${duration.getUTCMinutes()>1?'s':''}`;
}

export async function getLineColour(agency: string, line: string) {
  const res = await fetch(apiUrl('/api/getcolour?')+ new URLSearchParams({
    agency: agency,
    line: line
  }));
  let col = await res.text();
  if (col[0] == '#') {
    col = col.substring(1);
    col = lightenDarkenColor(col, 20);
    col = '#' + col;
  }
  return col
}

// https://stackoverflow.com/q/5560248
export function lightenDarkenColor(col, amt) {
  col = parseInt(col, 16);
  return (((col & 0x0000FF) + amt) | ((((col >> 8) & 0x00FF) + amt) << 8) | (((col >> 16) + amt) << 16)).toString(16);
}

// https://stackoverflow.com/a/11868398
export function getContrastYIQ(hexcolor){
  hexcolor = hexcolor.replace("#", "");
  const r = parseInt(hexcolor.substr(0,2),16);
  const g = parseInt(hexcolor.substr(2,2),16);
  const b = parseInt(hexcolor.substr(4,2),16);
  const yiq = ((r*299)+(g*587)+(b*114))/1000;
  return (yiq >= 128) ? 'black' : 'white';
}