import React from 'react';

interface ModalProps {
  message: string;
  onConfirm: () => void;
  onCancel: () => void;
}

const ModalDelete: React.FC<ModalProps> = ({ message, onConfirm, onCancel }) => {
  return (
    <div className="fixed inset-0 flex justify-center items-center bg-black bg-opacity-50">
      <div className="bg-white p-6 rounded shadow-md w-80">
        <p className="text-center text-lg mb-4">{message}</p>
        <div className="flex justify-between">
          <button
            onClick={onConfirm}
            className="bg-red-500 hover:bg-red-600 text-white py-2 px-4 rounded"
          >
            Удалить
          </button>
          <button
            onClick={onCancel}
            className="bg-gray-500 hover:bg-gray-600 text-white py-2 px-4 rounded"
          >
            Отменить
          </button>
        </div>
      </div>
    </div>
  );
};

export default ModalDelete;