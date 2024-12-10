const colorMap = new Map<string, string>();
const colors = ['processing', 'success', 'error', 'warning', 'magenta', 'red','volcano','orange','gold','lime','green','cyan','blue','geekblue','purple'];

export function getColor(text: string): string {
  if (!colorMap.has(text)) {
    const randomColor = colors[Math.floor(Math.random() * colors.length)];
    colorMap.set(text, randomColor);
  }
  return colorMap.get(text) || '';
}


export function formatTime(timeStr:string){
  // 将ISO 8601格式的时间字符串转换为指定格式
  const date = new Date(timeStr);
  
  // 获取年、月、日、时、分、秒
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  // 格式化为指定的字符串
  return `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`;
}