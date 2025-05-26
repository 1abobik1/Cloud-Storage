import React from 'react';




interface TranslateProps {
  type: string;
 
}


const TypeToTranslate: React.FC<TranslateProps> = ({ type }) => {



type TypeFileIcon = 'text' | 'photo' | 'video' |'unknown'|'home'  ;

const TypeMap: Record<TypeFileIcon, React.ReactNode> = {
    home: "Главная",
    text: "Документы",
    photo: "Фотография",
    video:  "Видео",
    unknown: "Прочие",
  };

const Type = (type.toLowerCase() as TypeFileIcon) in TypeMap ? (type.toLowerCase() as TypeFileIcon) : 'unknown';
  const Translate = TypeMap[Type];


    return (
         <span>
        {Translate}
        </span>
    );
};

export default TypeToTranslate;