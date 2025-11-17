import { Modal, Table } from "antd";
import type { ColumnsType } from "antd/es/table";
import { useEffect, useState } from "react";
import { groupMembers } from "../../api";

export interface MembersModalProps {
  isOpen: boolean
  chatroomId: number
  handleClose: () => void
  queryKey: string
}

interface User {
  id: number;
  username: string;
  nickName: string;
  headPic: string;
  email: string;
}

export function MembersModal(props: MembersModalProps) {

  const [members, setMembers] = useState<Array<User>>();

  const queryMembers = async () => {
    try {
      const res = await groupMembers(props.chatroomId);

      setMembers(res.data.map((item: User) => {
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

  useEffect(() => {
    queryMembers();
  }, [props.chatroomId, props.queryKey]);

  const columns: ColumnsType<User> = [
    {
      title: 'ID',
      dataIndex: 'id'
    },
    {
      title: '用户名',
      dataIndex: 'username'
    },
    {
      title: '昵称',
      dataIndex: 'nickName'
    },
    {
      title: '头像',
      dataIndex: 'headPic',
      render: (_, record) => (
        <div>
          <img src={record.headPic} width={50} height={50} />
        </div>
      )
    },
    {
      title: '邮箱',
      dataIndex: 'email'
    }
  ]

  return <Modal
    title="群聊成员"
    open={props.isOpen}
    onCancel={() => props.handleClose()}
    onOk={() => props.handleClose()}
    width={1000}
  >
    <Table columns={columns} dataSource={members} pagination={false} />
  </Modal>
}
