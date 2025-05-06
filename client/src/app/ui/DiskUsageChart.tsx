// components/DiskUsageChart.tsx
import { Pie } from "react-chartjs-2";
import { Chart as ChartJS, ArcElement, CategoryScale, Tooltip, Legend } from 'chart.js';

ChartJS.register(ArcElement, CategoryScale, Tooltip, Legend);

interface DiskUsageChartProps {
  fileCounts: Record<string, number>;
  totalSpace: number;
  totalUsedSpace: number;
}

const DiskUsageChart: React.FC<DiskUsageChartProps> = ({ fileCounts, totalSpace, totalUsedSpace }) => {
  const freeSpace = totalSpace - totalUsedSpace / 1024; 

  const chartData = {
    labels: ['Текст', 'Фото', 'Видео', 'Неизвестно', 'Свободное место'],
    datasets: [
      {
        data: [
          fileCounts.text,
          fileCounts.photo,
          fileCounts.video,
          fileCounts.unknown,
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
          <p className="text-xl">Свободно места: {freeSpace.toFixed(2)} GB</p>
        </div>
      </div>
    </div>
  );
};

export default DiskUsageChart;
