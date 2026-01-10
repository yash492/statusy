INSERT INTO component_groups(name, service_id)
VALUES 
    {{- range $i, $_ := .Values}}
        {if .eq .i le}
        (?, ?),
    {{- end}}