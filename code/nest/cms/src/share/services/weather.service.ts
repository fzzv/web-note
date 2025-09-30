import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import axios from 'axios';
import geoip from 'geoip-lite';

@Injectable()
export class WeatherService {
  constructor(private readonly configService: ConfigService) { }

  async getExternalIP() {
    try {
      const ipApiUrl = this.configService.get<string>('IP_API_URL')!;
      const response = await axios.get(ipApiUrl);
      return response.data.ip;
    } catch (error) {
      console.error('Error fetching external IP:', error);
      return 'N/A';
    }
  }

  async getWeather() {
    const ip = await this.getExternalIP();
    const geo = geoip.lookup(ip);
    const location = geo ? `${geo.city}, ${geo.country}` : 'Unknown';
    let weather = '无法获取天气信息';
    try {
      if (geo) {
        const apiKey = this.configService.get<string>('WEATHER_API_KEY');
        const weatherApiUrl = this.configService.get<string>('WEATHER_API_URL');
        const response = await axios.get(`${weatherApiUrl}?lang=zh&key=${apiKey}&q=${location}`);
        weather = `${response.data.current.temp_c}°C, ${response.data.current.condition.text}`;
      }
    } catch (error) {
      console.error('获取天气信息失败:', error.message);
    }
    return weather;
  }
}
