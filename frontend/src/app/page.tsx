"use client";

import { useEffect, useState } from "react";
import {
  Task,
  TaskType,
  getTasks,
  createTask,
  updateTaskCompleted,
  deleteTask,
} from "@/lib/api";

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [type, setType] = useState<TaskType>("adhoc");
  const [frequency, setFrequency] = useState("daily");

  async function loadTasks() {
    try {
      setLoading(true);
      const data = await getTasks();
      setTasks(data);
      setError(null);
    } catch {
      setError("اتصال به سرور برقرار نشد");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    // eslint-disable-next-line react-hooks/set-state-in-effect -- loadTasks sets state asynchronously after awaiting the API call, not synchronously within the effect body.
    loadTasks();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!title.trim()) return;

    await createTask({
      title,
      type,
      frequency: type === "routine" ? frequency : undefined,
    });
    setTitle("");
    loadTasks();
  }

  async function handleToggle(task: Task) {
    await updateTaskCompleted(task.id, !task.completed);
    loadTasks();
  }

  async function handleDelete(id: number) {
    await deleteTask(id);
    loadTasks();
  }

  return (
    <main className="max-w-2xl mx-auto p-6" dir="rtl">
      <h1 className="text-2xl font-bold mb-6">تسک‌های من</h1>

      <form onSubmit={handleSubmit} className="mb-8 space-y-3">
        <input
          type="text"
          placeholder="عنوان تسک"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          className="w-full border rounded px-3 py-2"
        />

        <div className="flex gap-3">
          <select
            value={type}
            onChange={(e) => setType(e.target.value as TaskType)}
            className="border rounded px-3 py-2"
          >
            <option value="adhoc">یک‌باره</option>
            <option value="routine">روتین</option>
          </select>

          {type === "routine" && (
            <select
              value={frequency}
              onChange={(e) => setFrequency(e.target.value)}
              className="border rounded px-3 py-2"
            >
              <option value="daily">روزانه</option>
              <option value="weekly">هفتگی</option>
            </select>
          )}

          <button
            type="submit"
            className="bg-black text-white px-4 py-2 rounded"
          >
            افزودن
          </button>
        </div>
      </form>

      {loading && <p>در حال بارگذاری...</p>}
      {error && <p className="text-red-600">{error}</p>}

      <ul className="space-y-2">
        {tasks.map((task) => (
          <li
            key={task.id}
            className="flex items-center justify-between border rounded px-4 py-2"
          >
            <div className="flex items-center gap-3">
              <input
                type="checkbox"
                checked={task.completed}
                onChange={() => handleToggle(task)}
              />
              <span className={task.completed ? "line-through text-gray-400" : ""}>
                {task.title}
              </span>
              <span className="text-xs text-gray-500">
                {task.type === "routine" ? `روتین (${task.frequency})` : "یک‌باره"}
              </span>
            </div>
            <button
              onClick={() => handleDelete(task.id)}
              className="text-red-600 text-sm"
            >
              حذف
            </button>
          </li>
        ))}
      </ul>

      {!loading && tasks.length === 0 && (
        <p className="text-gray-500">هنوز تسکی ثبت نشده.</p>
      )}
    </main>
  );
}
