
layui.use(['table', 'laydate', 'form'], function(){
    var table = layui.table
        ,form = layui.form
        ,util=layui.util
        ,laydate = layui.laydate;

    laydate.render({
        elem: '#start_time',
        type: 'datetime'
    });

    laydate.render({
        elem: '#end_time',
        type: 'datetime'
    });

    table.render({
        elem: '#roleList'
        ,url:'/admin/user/role_list'
        ,id: "roleListRender"
        ,toolbar: '#toolbarTpl'
        ,cellMinWidth: 80
        ,cols: [[
            {type: 'checkbox'}
            ,{field:'id', title:'ID', width:100, unresize: true, sort: true}
            ,{field:'parent_id', title:'父ID', width:100, unresize: true}
            ,{field:'role_name', title:'角色名', width:160}
            ,{field:'create_time', title:'创建时间', width:200, sort: true, templet: function(d) {return util.toDateString(d.create_time*1000); }}
            ,{field:'update_time', title:'更新时间', width:200, sort: true, templet: function(d) {return util.toDateString(d.update_time*1000); }}
            ,{field:'remark', title: '介绍', minWidth:120, sort: true}
            ,{field:'power', title: '权限', minWidth:120, sort: true,templet: '#powerTpl'}
            ,{field:'role_status', title:'操作', width:240, toolbar: '#actionTpl', unresize: true}
        ]]
        ,page: true
        // ,page:{ //支持传入 laypage 组件的所有参数（某些参数除外，如：jump/elem） - 详见文档
        //         layout: ['count', 'prev', 'page', 'next', 'skip', 'limit'] //自定义分页布局
        //             //,curr: 5 //设定初始在第 5 页
        //             ,groups: 5 //只显示 5 个连续页码
        //             ,first: true //显示首页
        //             ,last: true //显示尾页
        // }
    });

    form.on('submit(sreach)', function (data) {
        var loading = layer.load(1, {shade: [0.1, '#FF0000']});

        table.reload('roleListRender', {
            page: {
                curr: 1 //重新从第 1 页开始
            }
            ,where: data.field
        });
        layer.close(loading);
        return false
    })

    table.on('toolbar(roleListAction)', function(obj){
        var checkStatus = table.checkStatus(obj.config.id);
        var data = checkStatus.data;
        switch(obj.event){
            case 'addRole':
                x_admin_show('添加角色','/admin/user/role_edit');
                break;
            case 'startRoles':
                changeRoleStatusWithBatch(data,1);
                break;
            case 'stopRoles':
                changeRoleStatusWithBatch(data,0);
                break;
            case 'deleteRoles':
                deleteRoleWithBatch(data)
                break;
        };
    });

    //监听状态操作
    form.on('switch(statusSwitch)', function(obj){
        changeRoleStatus(this.value, obj)
    });

    table.on('tool(roleListAction)', function(obj){
        switch(obj.event){
            case 'delete':
                deleteRole(obj);
                break
            case 'edit':
                editRole(obj);
                break;
        }
    });


    function deleteRole(obj) {
        var data = obj.data;
        layer.confirm('确定删除吗', function(index){
            $.ajax({
                url: '/admin/user/role_delete',
                data: {"id": data.id,"name":data.role_name},
                type: "post",
                dataType: "json",
                success: function (ret) {
                    var message = ret.msg + ret.code;
                    if (ret.code === 0) {
                        message = ret.msg
                        obj.del();
                        layer.close(index);
                    }
                    layer.msg(message, {icon: 1, time: 1000}, function () {
                    });
                }
            });
        });
    }
    function deleteRoleWithBatch(data) {
        var roles = [];
        var roles_name = [];
        data.forEach(function(value,i) {
            roles.unshift(value.id)
            roles_name.unshift(value.role_name)
        });
        if (roles.length >0 ){
            layer.confirm('确定删除吗', function(index){
                $.ajax({
                    url: '/admin/user/role_delete',
                    data: {"id": roles,"name":roles_name},
                    type: "post",
                    dataType: "json",
                    traditional: true,
                    success: function (ret) {
                        var message = ret.msg + ret.code;
                        if (ret.code === 0) {
                            message = ret.msg
                            table.reload('roleListRender')
                        }
                        layer.msg(message, {icon: 1, time: 1000}, function () {
                        });
                    }
                });
            });
        }else{
            layer.msg("未选中目标", {icon: 1, time: 1000}, function () {});
        }
    }

    function editRole(obj){
        var data = obj.data;
        x_admin_show('添加角色','/admin/user/role_edit?role_name='+data.role_name+'&role_id='+data.id);
    }

    function changeRoleStatus(id,obj) {
        var status = 0
        if (obj.elem.checked) {
            status = 1
        }
        $.ajax({
            url: '/admin/user/role_status_change',
            data: {"id": id,"status":status},
            type: "post",
            dataType: "json",
            success: function (ret) {
                var message = ret.msg;
                if (ret.code < 0) {
                    message += ret.code;
                    obj.elem.checked = !obj.elem.checked;
                    form.render();
                }
                layer.msg(message, {icon: 1, time: 1000}, function () {});
            }
        });
    }

    function changeRoleStatusWithBatch(data, status) {
        var roles = [];
        data.forEach(function(value,i) {
            roles.push(value.id)
        });
        if (roles.length >0 ){
            $.ajax({
                url: '/admin/user/role_status_change',
                data: {"id": roles,"status":status},
                type: "post",
                dataType: "json",
                traditional: true,
                success: function (ret) {
                    var message = ret.msg + ret.code;
                    if (ret.code ===0) {
                        message = ret.msg
                        table.reload('roleListRender')
                    }
                    layer.msg(message, {icon: 1, time: 1000}, function () {});
                }
            });
        }else{
            layer.msg("未选中目标", {icon: 1, time: 1000}, function () {});
        }
    }
});
