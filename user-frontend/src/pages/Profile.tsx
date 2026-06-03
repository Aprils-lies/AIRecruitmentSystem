import { useState, useEffect } from 'react';
import { getProfile, updateProfile } from '../api/user';
import { normalizeProfile } from '../utils/nullsafe';
import type { UpdateProfileRequest } from '../types';
import Loading from '../components/Loading';

interface ProfileProps {
  onUpdateSuccess: () => void;
}

export default function Profile({ onUpdateSuccess }: ProfileProps) {
  const [loading, setLoading] = useState(false);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  
  const [profileData, setProfileData] = useState({
    real_name: '',
    phone: '',
    education: '',
    school: '',
    experience: '',
    skills: '',
  });
  
  const [formData, setFormData] = useState({
    real_name: '',
    phone: '',
    education: '',
    school: '',
    experience: '',
    skills: '',
  });

  useEffect(() => {
    setLoading(true);
    getProfile().then((res) => {
      const p = normalizeProfile(res.data);
      setProfileData({
        real_name: p.real_name,
        phone: p.phone,
        education: p.education,
        school: p.school,
        experience: p.experience,
        skills: p.skills,
      });
      setFormData({
        real_name: p.real_name,
        phone: p.phone,
        education: p.education,
        school: p.school,
        experience: p.experience,
        skills: p.skills,
      });
    }).catch((e) => {
      setError(e instanceof Error ? e.message : '加载个人资料失败');
    }).finally(() => {
      setLoading(false);
    });
  }, []);

  const handleChange = (field: keyof typeof formData, value: string) => {
    setFormData(prev => ({ ...prev, [field]: value }));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    setError(null);

    try {
      const updateData: UpdateProfileRequest = {};
      if (formData.real_name) updateData.real_name = formData.real_name;
      if (formData.phone) updateData.phone = formData.phone;
      if (formData.education) updateData.education = formData.education;
      if (formData.school) updateData.school = formData.school;
      if (formData.experience) updateData.experience = formData.experience;
      if (formData.skills) updateData.skills = formData.skills;

      await updateProfile(updateData);
      
      const res = await getProfile();
      const p = normalizeProfile(res.data);
      setProfileData({
        real_name: p.real_name,
        phone: p.phone,
        education: p.education,
        school: p.school,
        experience: p.experience,
        skills: p.skills,
      });
      setFormData({
        real_name: p.real_name,
        phone: p.phone,
        education: p.education,
        school: p.school,
        experience: p.experience,
        skills: p.skills,
      });
      
      setIsEditing(false);
      onUpdateSuccess();
    } catch (e) {
      setError(e instanceof Error ? e.message : '更新失败');
    } finally {
      setSaving(false);
    }
  };

  const handleCancel = () => {
    setFormData({ ...profileData });
    setIsEditing(false);
  };

  if (loading) {
    return <Loading text="加载个人资料..." />;
  }

  if (error) {
    return (
      <div className="bg-white rounded-lg shadow-md p-6">
        <p className="text-red-500">{error}</p>
      </div>
    );
  }

  const educationOptions = [
    '高中', '中专', '大专', '本科', '硕士', '博士', '其他'
  ];

  return (
    <div className="max-w-2xl mx-auto">
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">个人资料</h1>
        {!isEditing && (
          <button
            onClick={() => setIsEditing(true)}
            className="px-4 py-2 bg-primary-500 text-white rounded-lg hover:bg-primary-600"
          >
            编辑资料
          </button>
        )}
      </div>

      <div className="bg-white rounded-lg shadow-md p-6">
        {!isEditing ? (
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <p className="text-sm text-gray-500 mb-1">真实姓名</p>
                <p className="text-lg font-medium text-gray-900">{profileData.real_name || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 mb-1">联系电话</p>
                <p className="text-lg font-medium text-gray-900">{profileData.phone || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 mb-1">最高学历</p>
                <p className="text-lg font-medium text-gray-900">{profileData.education || '-'}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500 mb-1">毕业院校</p>
                <p className="text-lg font-medium text-gray-900">{profileData.school || '-'}</p>
              </div>
            </div>
            
            <div>
              <p className="text-sm text-gray-500 mb-1">工作/项目经历</p>
              <p className="text-gray-900 whitespace-pre-line">{profileData.experience || '-'}</p>
            </div>
            
            <div>
              <p className="text-sm text-gray-500 mb-1">核心技能标签</p>
              <div className="flex flex-wrap gap-2">
                {profileData.skills ? (
                  profileData.skills.split(',').map((skill, index) => (
                    <span
                      key={index}
                      className="px-3 py-1 bg-primary-100 text-primary-700 rounded-full text-sm"
                    >
                      {skill.trim()}
                    </span>
                  ))
                ) : (
                  <span className="text-gray-400">-</span>
                )}
              </div>
            </div>
          </div>
        ) : (
          <>
            {error && (
              <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg">
                {error}
              </div>
            )}

            <form onSubmit={handleSubmit}>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">真实姓名 *</label>
                  <input
                    type="text"
                    value={formData.real_name}
                    onChange={(e) => handleChange('real_name', e.target.value)}
                    placeholder="请输入真实姓名"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">联系电话 *</label>
                  <input
                    type="tel"
                    value={formData.phone}
                    onChange={(e) => handleChange('phone', e.target.value)}
                    placeholder="请输入联系电话"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">最高学历 *</label>
                  <select
                    value={formData.education}
                    onChange={(e) => handleChange('education', e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  >
                    <option value="">请选择学历</option>
                    {educationOptions.map((edu) => (
                      <option key={edu} value={edu}>{edu}</option>
                    ))}
                  </select>
                </div>

                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">毕业院校 *</label>
                  <input
                    type="text"
                    value={formData.school}
                    onChange={(e) => handleChange('school', e.target.value)}
                    placeholder="请输入毕业院校"
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                  />
                </div>
              </div>

              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">工作/项目经历 *</label>
                <textarea
                  value={formData.experience}
                  onChange={(e) => handleChange('experience', e.target.value)}
                  placeholder="请详细描述您的工作或项目经历"
                  rows={4}
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 resize-none"
                />
              </div>

              <div className="mb-6">
                <label className="block text-sm font-medium text-gray-700 mb-2">核心技能标签 *</label>
                <input
                  type="text"
                  value={formData.skills}
                  onChange={(e) => handleChange('skills', e.target.value)}
                  placeholder="请输入核心技能，用逗号分隔"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500"
                />
              </div>

              <div className="flex gap-3">
                <button
                  type="button"
                  onClick={handleCancel}
                  className="flex-1 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50"
                >
                  取消
                </button>
                <button
                  type="submit"
                  disabled={saving}
                  className="flex-1 py-3 bg-primary-500 text-white rounded-lg hover:bg-primary-600 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {saving ? '保存中...' : '保存'}
                </button>
              </div>
            </form>
          </>
        )}
      </div>
    </div>
  );
}