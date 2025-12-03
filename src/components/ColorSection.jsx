import { useState, useEffect } from 'react';
import { MdPalette } from 'react-icons/md';
import { HexColorPicker } from 'react-colorful';

const ColorSection = ({ publishMessage }) => {
  const [color, setColor] = useState(localStorage.getItem('custom-color') || '#1f2937'); // default به dark

  useEffect(() => {
    localStorage.setItem('custom-color', color);
    publishMessage('color/set', { color });
  }, [color, publishMessage]);

  return (
    <div className="bg-gray-800 p-6 rounded-xl shadow-lg border border-gray-600 transition-all duration-300 hover:shadow-xl animate-fade-in">
      <h2 className="text-xl font-bold mb-6 text-center text-gray-100 flex items-center justify-center gap-2">
        <MdPalette size={24} className="text-color" /> انتخاب رنگ
      </h2>
      
      <div className="mb-6 flex justify-center items-center flex-col gap-2">
        <div className='flex w-full justify-between'>
        <label className="block text-color font-bold mb-3 text-gray-300 text-right">
          رنگ پس‌زمینه سفارشی:
        </label>
        <span className="text-color font-bold text-gray-300">{color}</span>
        </div>
        <div className='flex rounded-xl shadow-xl border border-gray-600 p-6'>
          <HexColorPicker
              color={color}
              onChange={setColor}
              className="react-colorful-wrapper" // class برای CSS override
            />
        </div>
        
      </div>
    </div>
  );
};

export default ColorSection;