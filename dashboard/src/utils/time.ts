/**
 * Formats uptime seconds into a human-readable string
 * @param seconds - Total uptime in seconds
 * @param useKorean - Whether to use Korean labels (default: false)
 * @returns Formatted string like "2d 5h 30m 15s" or "2일 5시간 30분 15초"
 */
export const formatUptime = (seconds: number, useKorean: boolean = false): string => {
  if (seconds < 0) return useKorean ? "0초" : "0s";

  const days = Math.floor(seconds / 86400);
  const hours = Math.floor((seconds % 86400) / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  const secs = Math.floor(seconds % 60);

  const parts: string[] = [];

  if (useKorean) {
    if (days > 0) parts.push(`${days}일`);
    if (hours > 0) parts.push(`${hours}시간`);
    if (minutes > 0) parts.push(`${minutes}분`);
    if (secs > 0 || parts.length === 0) parts.push(`${secs}초`);
  } else {
    if (days > 0) parts.push(`${days}d`);
    if (hours > 0) parts.push(`${hours}h`);
    if (minutes > 0) parts.push(`${minutes}m`);
    if (secs > 0 || parts.length === 0) parts.push(`${secs}s`);
  }

  return parts.join(" ");
};
