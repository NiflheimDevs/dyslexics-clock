import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import {
  getAlarms,
  createAlarm,
  updateAlarm,
  deleteAlarm,
  getDeviceColor,
  updateDeviceColor,
  getDeviceVolume,
  updateDeviceVolume,
  BezanBekob,
  GetBrightness,
  SnoozeWAlarm,
  StopAlarm,
  UpdateBrightness,
} from "../api/services";
import { toast } from "react-toastify";

export const useAlarms = () => {
  return useQuery({
    queryKey: ["alarms"],
    queryFn: getAlarms,
    select: (data) => data.data,
  });
};

export const useCreateAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: createAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("آلارم با موفقیت ایجاد شد");
    },
    onError: () => toast.error("خطا در ایجاد آلارم"),
  });
};

export const useUpdateAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, data }) => updateAlarm(id, data),
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("آلارم با موفقیت بروز شد");
    },
  });
};

export const useUpdateVolume = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateDeviceVolume,
    onSuccess: () => {
      queryClient.invalidateQueries(["deviceVolume"]);
      toast.success("میزان صدای دستگاه با موفقیت بروز شد");
    },
  });
};

export const useDeleteAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: deleteAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("آلارم با موفقیت حذف شد");
    },
  });
};

export const useDeviceColor = () => {
  return useQuery({
    queryKey: ["deviceColor"],
    queryFn: getDeviceColor,
    select: (data) => data.data.color,
  });
};

export const useDeviceVolume = () => {
  return useQuery({
    queryKey: ["deviceVolume"],
    queryFn: getDeviceVolume,
    select: (data) => data.data.volume,
  });
};

export const useUpdateColor = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: updateDeviceColor,
    onSuccess: () => {
      queryClient.invalidateQueries(["deviceColor"]);
      toast.success("رنگ با موفقیت بروز شد");
    },
  });
};

export const useDeviceBrightness = () => {
  return useQuery({
    queryKey: ["deviceBrightness"],
    queryFn: GetBrightness,
    select: (data) => data.data.brightness,
  });
};

export const useUpdateBrightness = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: UpdateBrightness,
    onSuccess: () => {
      queryClient.invalidateQueries(["deviceBrightness"]);
      toast.success("روشنایی با موفقیت بروز شد");
    },
  });
};

export const useSnoozeAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: SnoozeWAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("آلارم به تعویق افتاد");
    },
    onError: () => toast.error("خطا در به تعویق انداختن آلارم"),
  });
};

export const useStopAlarm = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: StopAlarm,
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("آلارم متوقف شد");
    },
    onError: () => toast.error("خطا در متوقف کردن آلارم"),
  });
};

export const useBezanBekob = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: BezanBekob,
    onSuccess: () => {
      queryClient.invalidateQueries(["alarms"]);
      toast.success("بزن بکوب با موفقیت انجام شد");
    },
    onError: () => toast.error("خطا در انجام بزن بکوب"),
  });
};
