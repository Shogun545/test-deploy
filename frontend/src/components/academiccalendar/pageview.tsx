"use client";
import React, { useState, useEffect } from "react";
import { Button, Modal, Form, Input, DatePicker, Select, message, Spin, Popconfirm } from "antd";
import { LeftOutlined, RightOutlined, ClockCircleOutlined, EditOutlined, DeleteOutlined, PlusOutlined } from "@ant-design/icons";
import dayjs, { Dayjs } from "dayjs";
import isBetween from "dayjs/plugin/isBetween";
import "dayjs/locale/th";

// Setup dayjs
dayjs.extend(isBetween);
dayjs.locale("th");

const { RangePicker } = DatePicker;

interface CalendarEvent {
  id: string | number;
  startDate: string;
  endDate: string;
  startTime: string;
  endTime: string;
  EventName: string;
  EventType: "exam" | "holiday" | "activity";
  isFixed: boolean;
}

interface CalendarViewProps {
  canEdit: boolean;
}

export default function CalendarView({ canEdit }: CalendarViewProps) {
  const [currentDate, setCurrentDate] = useState(dayjs());
  const [selectedDate, setSelectedDate] = useState<string>(dayjs().format("YYYY-MM-DD"));
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [editingEvent, setEditingEvent] = useState<CalendarEvent | null>(null);
  const [events, setEvents] = useState<CalendarEvent[]>([]);
  const [form] = Form.useForm();
  const [messageApi, contextHolder] = message.useMessage();

  const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

  const disabledDate = (current: Dayjs) => {
    return current && current < dayjs().startOf('day');
  };

  const fetchEvents = async (year: number) => {
    setLoading(true);
    try {
      const res = await fetch(`${API_URL}/api/holidays?year=${year}`);
      if (!res.ok) throw new Error("Network response was not ok");

      const data = await res.json();

      const mappedEvents: CalendarEvent[] = data.map((item: any) => {
        const sDate = item.start_date || item.date;
        const eDate = item.end_date || item.date;
        const sTime = item.start_time || "00:00";
        const eTime = item.end_time || "23:59";

        const rawType = (item.types && item.types[0]) ? item.types[0] : 'exam';
        let finalType: "exam" | "holiday" | "activity" = 'exam';
        let fixed = false;

        if (rawType === 'Public Holiday' || rawType === 'holiday') {
          finalType = 'holiday';
          fixed = rawType === 'Public Holiday';
        } else if (rawType === 'activity') {
          finalType = 'activity';
        } else if (rawType === 'exam') {
          finalType = 'exam';
        }

        return {
          id: item.id || `evt-${Math.random()}`,
          startDate: sDate,
          endDate: eDate,
          startTime: sTime,
          endTime: eTime,
          EventName: item.localName || item.event_name || item.name,
          EventType: finalType,
          isFixed: fixed,
        };
      });

      setEvents(mappedEvents);
    } catch (error) {
      console.error(error);
      messageApi.error("ไม่สามารถดึงข้อมูลปฏิทินได้");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents(currentDate.year());
  }, [currentDate.year()]);

  const openAddModal = (dateStr?: string) => {
    if (!canEdit) return;
    setEditingEvent(null);
    form.resetFields();
    const targetDate = dateStr ? dayjs(dateStr) : dayjs(selectedDate);
    const defaultStart = targetDate.hour(0).minute(0);
    const defaultEnd = targetDate.hour(23).minute(59);
    form.setFieldsValue({ dateRange: [defaultStart, defaultEnd], type: undefined });
    setIsModalOpen(true);
  };

  const openEditModal = (event: CalendarEvent) => {
    if (!canEdit) return;
    setEditingEvent(event);
    const startObj = dayjs(`${event.startDate}T${event.startTime}`);
    const endObj = dayjs(`${event.endDate}T${event.endTime}`);
    form.setFieldsValue({ title: event.EventName, dateRange: [startObj, endObj], type: event.EventType });
    setIsModalOpen(true);
  };

  const handleSaveEvent = async (values: any) => {
    if (!canEdit) return;
    if (!values.dateRange || values.dateRange.length < 2) {
      messageApi.error("กรุณาระบุวันและเวลาเริ่มต้น-สิ้นสุด");
      return;
    }
    const startObj = values.dateRange[0];
    const endObj = values.dateRange[1];
    const payload = {
      title: values.title, type: values.type,
      start_date: startObj.format("YYYY-MM-DD"), end_date: endObj.format("YYYY-MM-DD"),
      start_time: startObj.format("HH:mm"), end_time: endObj.format("HH:mm"),
    };
    const token = localStorage.getItem("token");
    if (!token) { messageApi.error("กรุณาเข้าสู่ระบบก่อนทำรายการ"); return; }

    try {
      const method = editingEvent ? 'PUT' : 'POST';
      const url = editingEvent ? `${API_URL}/api/events/${editingEvent.id}` : `${API_URL}/api/events`;
      const res = await fetch(url, { method, headers: { 'Content-Type': 'application/json', 'Authorization': `Bearer ${token}` }, body: JSON.stringify(payload) });
      if (!res.ok) throw new Error(await res.text());
      messageApi.success("บันทึกข้อมูลเรียบร้อย");
      setIsModalOpen(false); setEditingEvent(null); form.resetFields();
      fetchEvents(currentDate.year());
    } catch (err: any) { console.error(err); messageApi.error(`บันทึกข้อมูลไม่สำเร็จ: ${err.message}`); }
  };

  const handleDeleteEvent = async (id: string | number) => {
    if (!canEdit) return;
    const token = localStorage.getItem("token");
    if (!token) return;
    try {
      const res = await fetch(`${API_URL}/api/events/${id}`, { method: 'DELETE', headers: { 'Authorization': `Bearer ${token}` } });
      if (!res.ok) throw new Error("Delete failed");
      messageApi.success("ลบกิจกรรมเรียบร้อย");
      fetchEvents(currentDate.year());
    } catch (err) { messageApi.error("ลบข้อมูลไม่สำเร็จ"); }
  };

  const getEventTag = (type: string) => {
    switch (type) {
      case 'exam': return { label: 'การสอบ', className: 'bg-yellow-100 text-yellow-700', iconColor: 'text-yellow-500' };
      case 'activity': return { label: 'กิจกรรม', className: 'bg-orange-100 text-orange-600', iconColor: 'text-orange-500' };
      case 'holiday': return { label: 'วันหยุด', className: 'bg-red-100 text-red-600', iconColor: 'text-red-500' };
      default: return { label: 'อื่นๆ', className: 'bg-gray-100 text-gray-600', iconColor: 'text-gray-500' };
    }
  };

  const activeEvents = events.filter((e) => {
    return dayjs(selectedDate).isBetween(e.startDate, e.endDate, 'day', '[]');
  });

  const daysInMonth = currentDate.daysInMonth();
  const startDay = currentDate.startOf("month").day();
  const calendarCells = [...Array(startDay).fill(null), ...Array(daysInMonth).keys().map((i) => i + 1)];

  return (
    <div className="h-[calc(100vh-120px)] flex flex-col rounded-3xl to-white">
      {contextHolder}

      <div className="flex justify-end items-center mb-2"></div>
      
      <div className="flex gap-8 h-full items-stretch overflow-hidden"> 

        {/* LEFT SIDE: CALENDAR */}
        <div className="flex flex-[2] flex-col h-full">
          <div className="flex items-center gap-3 mb-6 shrink-0">
            <div className="w-1 h-8 bg-gradient-to-b from-orange-500 to-orange-600 rounded-full"></div>
            <h2 className="text-3xl font-bold bg-gradient-to-r from-orange-600 to-orange-500 bg-clip-text text-transparent">
              {currentDate.format("MMMM YYYY")}
            </h2>
            {loading && <Spin size="small" className="ml-2" />}
            <div className="flex bg-gradient-to-r from-orange-500 to-orange-600 rounded-xl p-1.5 gap-1 ml-auto shadow-sm ">
              <Button type="text"  icon={<LeftOutlined />} onClick={() => setCurrentDate(currentDate.add(-1, 'month'))} className="!text-white hover:!bg-white hover:!text-orange-600 hover:shadow-sm !w-12 !h-12 !text-xl flex items-center justify-center" />
              <Button type="text"  icon={<RightOutlined />} onClick={() => setCurrentDate(currentDate.add(1, 'month'))} className="!text-white hover:!bg-white hover:!text-orange-600 hover:shadow-sm !w-12 !h-12 !text-xl flex items-center justify-center" />
            </div>
          </div>
          
          <div className="flex-1 flex flex-col rounded-2xl overflow-hidden shadow-sm text-sm border border-gray-100">
            <div className="grid grid-cols-7 bg-gradient-to-r from-orange-500 to-orange-600 shrink-0">
              {['อา', 'จ', 'อ', 'พ', 'พฤ', 'ศ', 'ส'].map(day => (
                <div key={day} className="py-4 text-center font-bold text-white text-base">{day}</div>
              ))}
            </div>
            <div className="grid grid-cols-7 grid-rows-5 flex-1 bg-white">
              {calendarCells.map((day, idx) => {
                const dateObj = day ? currentDate.date(day) : null;
                const dateStr = dateObj ? dateObj.format("YYYY-MM-DD") : "";
                const dayEvents = dateObj ? events.filter(e => dateObj.isBetween(e.startDate, e.endDate, 'day', '[]')) : [];
                const isSelected = selectedDate === dateStr;
                const isPastDate = dateObj ? dateObj.isBefore(dayjs(), 'day') : false;
                const isToday = dateObj ? dateObj.isSame(dayjs(), 'day') : false;

                return (
                  <div key={idx}
                    onClick={() => {
                      if (isPastDate) return;
                      if (day) {
                        setSelectedDate(dateStr);
                        if (canEdit) openAddModal(dateStr);
                      }
                    }}
                    className={`border-r border-b border-orange-50 p-2 relative transition-all duration-200
                        ${!day ? 'bg-gray-50/50' : ( isPastDate ? 'bg-gray-50 cursor-not-allowed opacity-60' : 'bg-white cursor-pointer hover:bg-orange-50' )} 
                        ${isSelected ? 'ring-2 ring-orange-500 z-10 shadow-md bg-white' : ''}
                        ${isToday && !isSelected ? 'ring-2 ring-orange-300' : ''}`}>
                    {day && (
                      <>
                        <span className={`font-bold block text-right mb-1 text-sm
                            ${isSelected ? 'text-orange-600' : ( isToday ? 'text-orange-500' : ( isPastDate ? 'text-gray-400' : 'text-gray-700' ) )}`}>
                          {day}
                        </span>
                        <div className="flex flex-col gap-1 mt-1">
                          {dayEvents.slice(0, 3).map((ev, i) => {
                            const tagStyle = getEventTag(ev.EventType);
                            return (
                              <div key={i} onClick={(e) => { e.stopPropagation(); if (!isPastDate) setSelectedDate(dateStr); }}
                                className={`text-[10px] px-1.5 py-0.5 rounded-md truncate font-medium shadow-sm ${tagStyle.className}`} title={ev.EventName}>
                                {ev.EventName}
                              </div>
                            );
                          })}
                          {dayEvents.length > 3 && <div className="text-[9px] text-gray-400 text-center">+{dayEvents.length - 3} ..</div>}
                        </div>
                      </>
                    )}
                  </div>
                );
              })}
            </div>
          </div>
        </div>

        {/* RIGHT SIDE: DETAILS PANEL */}
        <div className="w-[360px] h-full bg-gradient-to-br from-orange-500 via-orange-600 to-orange-700 rounded-3xl p-6 text-white shadow-2xl flex flex-col relative overflow-hidden">
          <div className="absolute top-0 right-0 w-40 h-40 bg-white/10 rounded-full blur-3xl"></div>
          <div className="absolute bottom-0 left-0 w-32 h-32 bg-white/10 rounded-full blur-2xl"></div>

          <div className="relative z-10 mb-4 text-center border-b border-white/30 pb-4 shrink-0">
            <h3 className="text-2xl font-bold mb-1">{dayjs(selectedDate).format("D MMMM YYYY")}</h3>
            <p className="text-sm opacity-90 bg-white/20 inline-block px-4 py-1 rounded-full font-medium">
              มี {activeEvents.length} กิจกรรม
            </p>
          </div>

          <div className="flex-1 overflow-y-auto pr-2 flex flex-col gap-3 custom-scrollbar relative z-10">
            {activeEvents.length === 0 ? (
              <div className="flex flex-col items-center justify-center h-full text-center">
                <div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center mb-4">
                  <ClockCircleOutlined className="text-3xl text-white/60" />
                </div>
                <p className="text-white/70 text-sm">ไม่มีกิจกรรมในวันนี้</p>
              </div>
            ) : (
              activeEvents.map((event) => {
                const isEventPast = dayjs(event.endDate).isBefore(dayjs(), 'day');
                const tagInfo = getEventTag(event.EventType);

                return (
                  <div key={event.id} className="bg-white rounded-xl p-4 text-gray-800 shadow-lg hover:shadow-xl transition-all duration-200">
                    <div className="flex justify-between items-start mb-2">
                      <h4 className="font-bold text-base leading-tight w-[70%]">{event.EventName}</h4>
                      <span className={`${tagInfo.className} text-[10px] px-2 py-0.5 rounded-full whitespace-nowrap font-semibold shadow-sm`}>
                        {tagInfo.label}
                      </span>
                    </div>

                    <div className="text-xs text-gray-600 mb-2">
                        {/* 🔥 แก้ไข: แสดง Start/End แบบเต็มเสมอ ไม่เช็คเงื่อนไขแล้ว */}
                        <div className="bg-gray-50 p-2 rounded-lg space-y-1">
                          <div className="flex justify-between">
                            <span className="font-medium text-gray-500">เริ่ม:</span>
                            <span className="font-semibold">{dayjs(event.startDate).format("D MMM")} {event.startTime}</span>
                          </div>
                          <div className="flex justify-between">
                            <span className="font-medium text-gray-500">สิ้นสุด:</span>
                            <span className="font-semibold">{dayjs(event.endDate).format("D MMM")} {event.endTime}</span>
                          </div>
                        </div>
                    </div>

                    {canEdit && !event.isFixed && !isEventPast && (
                      <div className="flex gap-2 mt-3 pt-2 border-t border-gray-100">
                        <Button size="small" className="bg-yellow-400 text-black border-none flex-1 font-medium hover:bg-yellow-500" icon={<EditOutlined />} onClick={() => openEditModal(event)}>แก้ไข</Button>
                        <Popconfirm title="ยืนยันการลบ" onConfirm={() => handleDeleteEvent(event.id)} okText="ลบ" cancelText="ยกเลิก" okButtonProps={{ danger: true }}>
                          <Button size="small" type="primary" danger className="flex-1 font-medium" icon={<DeleteOutlined />}>ลบ</Button>
                        </Popconfirm>
                      </div>
                    )}
                  </div>
                );
              }))}
          </div>
        </div>

        {/* MODAL */}
        <Modal
          title={<div className="flex items-center gap-3"><div className="w-1 h-6 bg-orange-500 rounded-full"></div><span className="text-xl font-bold">{editingEvent ? "แก้ไขกิจกรรม" : "เพิ่มกิจกรรมใหม่"}</span></div>}
          open={isModalOpen}
          onCancel={() => { setIsModalOpen(false); setEditingEvent(null); form.resetFields(); }}
          footer={null} centered width={500}
        >
          <Form form={form} layout="vertical" className="mt-6" onFinish={handleSaveEvent}>
            <Form.Item name="title" label="ชื่อกิจกรรม" rules={[{ required: true, message: 'โปรดระบุชื่อกิจกรรม' }]}>
              <Input size="middle" placeholder="เช่น สอบกลางภาค" />
            </Form.Item>
            <Form.Item name="dateRange" label="ระยะเวลา" rules={[{ required: true, message: 'โปรดเลือกช่วงเวลา' }]}>
              <RangePicker size="middle" showTime={{ format: 'HH:mm' }} format="YYYY-MM-DD HH:mm" className="w-full" placeholder={['วันเวลาเริ่ม', 'วันเวลาสิ้นสุด']} disabledDate={disabledDate} />
            </Form.Item>
            <Form.Item name="type" label="ประเภท" initialValue="exam" rules={[{ required: true, message: 'กรุณาระบุประเภทกิจกรรม' }]} >
              <Select size="middle" showSearch allowClear placeholder="เลือกประเภทกิจกรรม" options={[{ label: 'สอบ', value: 'exam' }, { label: 'กิจกรรม', value: 'activity' }, { label: 'วันหยุด', value: 'holiday' }]} />
            </Form.Item>
            <div className="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200">
              <Button onClick={() => setIsModalOpen(false)}>ยกเลิก</Button>
              <Button type="primary" htmlType="submit" className="bg-orange-500">{editingEvent ? "บันทึกการแก้ไข" : "บันทึก"}</Button>
            </div>
          </Form>
        </Modal>

      </div>
      <style jsx global>{`
        .custom-scrollbar::-webkit-scrollbar { width: 4px; }
        .custom-scrollbar::-webkit-scrollbar-track { background: rgba(255, 255, 255, 0.1); }
        .custom-scrollbar::-webkit-scrollbar-thumb { background: rgba(255, 255, 255, 0.3); border-radius: 10px; }
      `}</style>
    </div>
  );
}