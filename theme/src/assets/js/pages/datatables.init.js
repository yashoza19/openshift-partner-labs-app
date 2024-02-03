$(document).ready(function() {
    $.fn.dataTable.moment('MM/DD/YYYY');
    $.fn.dataTable.moment('MM/DD/YYYY HH:mm:ss');

    $('#requests-table').DataTable({
        columns: [
            { orderable: false },
            { orderable: false },
            null,
            null,
            { orderable: false },
            null,
            { orderable: false },
        ],
        order: [[ 2, "desc" ]],
    });

    $('#archive-requests-table').DataTable({
        columns: [
            { orderable: false },
            { orderable: false },
            null,
        ],
        order: [[ 2, "desc" ]],
    });

    $('#clusterlogins').DataTable({
        columns: [
            { orderable: false },
            { orderable: false },
            { orderable: false },
        ],
        order: [[ 2, "desc" ]],
        searching: false,
        lengthChange: false,
        pageLength: 5,
        info: false,
    });
    $('#extensions').DataTable({
        columns: [
            { orderable: false },
            { orderable: false },
            { orderable: false },
        ],
        order: [[ 1, "desc" ]],
        searching: false,
        lengthChange: false,
        pageLength: 10,
        info: false,
    });
});