// components/DiskUsageChart.tsx
import { Pie } from "react-chartjs-2";
import { Chart as ChartJS, ArcElement, CategoryScale, Tooltip, Legend } from 'chart.js';

ChartJS.register(ArcElement, CategoryScale, Tooltip, Legend);
import { useEffect, useState } from 'react';
import { jwtDecode } from 'jwt-decode';

interface UsageResponse {
  current_used_gb: number;
  current_used_mb: number;
  current_used_kb: number;
  storage_limit_gb: number;
  plan_name: string;
}

interface JwtPayload {
  user_id: number;
}
interface DiskUsageChartProps {
  fileCounts: Record<string, number>;
  totalUsedSpace: number;
}


const DiskUsageChart: React.FC<DiskUsageChartProps> = ({ fileCounts, totalUsedSpace }) => {
  

  const [usage, setUsage] = useState<UsageResponse | null>(null);
    const[loading, setLoading] = useState(true);

  

  useEffect(() => {
    const token = localStorage.getItem('token');
    if (!token) return;

    const decoded: JwtPayload = jwtDecode(token);
    const userId = decoded.user_id;

    fetch(`http://localhost:8081/user/${userId}/usage`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then(res => {
        if (!res.ok) throw new Error(`Ошибка: ${res.status}`);
        return res.json();
      })
      .then(data => {
        console.log('Ответ usage:', data);
        setUsage(data);
      })
      .catch(err => {
        console.error('Ошибка загрузки использования:', err);
      })
      .finally(() => {
        setLoading(false);
      });
  }, []);
  if (loading || !usage) return <div className="text-gray-600">Загрузка...</div>;
  const {current_used_gb,current_used_mb,storage_limit_gb } = usage;
  const ostatok =storage_limit_gb - (current_used_gb + current_used_mb / 1024)
  const freeSpace = storage_limit_gb - totalUsedSpace / 1024; 

  const chartData = {
    labels: ['Текст', 'Фото', 'Видео', 'Прочие', 'Свободное место'],
    datasets: [
      {
        data: [
          fileCounts.text,
          fileCounts.photo,
          fileCounts.video,
          fileCounts.other,
          freeSpace,
          
        ],
        backgroundColor: ['#4c6ef5', '#d8e2dc', '#e06c75', '#f5a623', '#a0aec0'], // Цвет для свободного места
        borderColor: '#fff',
        borderWidth: 1,
      }
    ]
  };

  return (
    <div className="mt-6">
      <h3 className="text-lg mb-4">🗂 Статистика использования</h3>
      <div className="flex items-center justify-between">
        <div className="w-1/2">
          <Pie data={chartData} options={{ responsive: true }} />
        </div>
        <div className="w-1/2 pl-6">
          <p className="text-xl">Свободно места: {ostatok.toFixed(2)} GB</p>
        </div>
      </div>
    </div>
  );
};

export default DiskUsageChart;
