import { Button, Form, Input, Table } from "antd";
import { useEffect, useMemo, useState } from "react";
import './index.css';
import type { ColumnsType } from "antd/es/table";
import { useForm } from "antd/es/form/Form";
import { chatroomList } from "../../api";

interface SearchGroup {
  name: string;
}

interface GroupSearchResult {
  id: number;
  name: string;
  createTime: Date;
}

export function Group() {
  const [groupResult, setGroupResult] = useState<Array<GroupSearchResult>>([]);

  const columns: ColumnsType<GroupSearchResult> = useMemo(() => [
    {
      title: '名称',
      dataIndex: 'name'
    },
    {
      title: '创建时间',
      dataIndex: 'createTime'
    },
    {
      title: '操作',
      render: () => (
        <div>
          <a href="#">聊天</a>
        </div>
      )
    }
  ], []);

  const searchGroup = async (values: SearchGroup) => {
    try {
      const res = await chatroomList(values.name || '');

      setGroupResult(res.data.map((item: GroupSearchResult) => {
        return {
          ...item,
          key: item.id
        }
      }));
    } catch (e: unknown) {
      if (e instanceof Error) {
        console.error(e.message);
      } else {
        console.error('Unknown error:', e);
      }
    }
  };

  const [form] = useForm();

  useEffect(() => {
    searchGroup({
      name: form.getFieldValue('name')
    });
  }, []);


  return <div id="group-container">
    <div className="group-form">
      <Form
        form={form}
        onFinish={searchGroup}
        name="search"
        layout='inline'
        colon={false}
      >
        <Form.Item label="名称" name="name">
          <Input />
        </Form.Item>

        <Form.Item label=" ">
          <Button type="primary" htmlType="submit">
            搜索
          </Button>
        </Form.Item>
      </Form>
    </div>
    <div className="group-table">
      <Table columns={columns} dataSource={groupResult} style={{ width: '1000px' }} />
    </div>
  </div>
}
