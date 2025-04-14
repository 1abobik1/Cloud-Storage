import {FileData} from '../FileData'

export interface CloudResponse {
    file_data: FileData[];
    message: string;
    status: number;
}

