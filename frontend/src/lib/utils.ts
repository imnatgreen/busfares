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
  col = lightenDarkenColor(col, 20);
  return col
}

// https://stackoverflow.com/q/5560248
export function lightenDarkenColor(col, amt) {
  let usePound = false;
  if ( col[0] == "#" ) {
      col = col.slice(1);
      usePound = true;
  }

  const num = parseInt(col,16);

  let r = (num >> 16) + amt;

  if ( r > 255 ) r = 255;
  else if  (r < 0) r = 0;

  let b = ((num >> 8) & 0x00FF) + amt;

  if ( b > 255 ) b = 255;
  else if  (b < 0) b = 0;
  
  let g = (num & 0x0000FF) + amt;

  if ( g > 255 ) g = 255;
  else if  ( g < 0 ) g = 0;

  return (usePound?"#":"") + (g | (b << 8) | (r << 16)).toString(16);  
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