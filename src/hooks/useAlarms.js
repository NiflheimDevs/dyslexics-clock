import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import {
  getAlarms,
  createAlarm,
  updateAlarm,
  deleteAlarm,
  getDeviceColor,
  updateDeviceColor,
} from '../api/services';
import { toast } from 'react-toastify';

export const useAlarms = () => {
  return useQuery({
    queryKey: ['alarms'],
    queryFn: getAlarms,
    select: (data) => data.data,
  });
};

export const useCreateAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(['alarms']);
      toast.success('آلارم با موفقیت ایجاد شد');
    },
    onError: () => toast.error('خطا در ایجاد آلارم'),
  });
};

export const useUpdateAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }) => updateAlarm(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries(['alarms']);
      toast.success('آلارم با موفقیت بروز شد');
    },
  });
};

export const useDeleteAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(['alarms']);
      toast.success('آلارم با موفقیت حذف شد');
    },
  });
};

export const useDeviceColor = () => {
  return useQuery({
    queryKey: ['deviceColor'],
    queryFn: getDeviceColor,
    select: (data) => data.data.color,
  });
};

export const useUpdateColor = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateDeviceColor,
    onSuccess: () => {
      queryClient.invalidateQueries(['deviceColor']);
      toast.success('رنگ با موفقیت بروز شد');
    },
  });
};