(module

  (func $luminance128v (param $fcolor v128) (result f32)
    (local $v v128)

    ;; square
    local.get $fcolor
    local.get $fcolor
    f32x4.mul

    ;; R: 0.2126 * (1.0/255.0)^2
    ;; G: 0.7152 * (1.0/255.0)^2
    ;; B: 0.0722 * (1.0/255.0)^2
    v128.const f32x4 3.269512e-6 1.099885e-5 1.110342e-6 0.0
    f32x4.mul
    local.set $v

    local.get $v
    f32x4.extract_lane 0

    local.get $v
    f32x4.extract_lane 1

    local.get $v
    f32x4.extract_lane 2

    f32.add
    f32.add
  )

  (func $luminance128f (param $fcolor v128) (result f32)
    (local $v v128)

    ;; square
    local.get $fcolor
    local.get $fcolor
    f32x4.mul

    ;; R: 0.2126
    ;; G: 0.7152
    ;; B: 0.0722
    v128.const f32x4 0.2126 0.7152 0.0722 0
    f32x4.mul
    local.set $v

    local.get $v
    f32x4.extract_lane 0

    local.get $v
    f32x4.extract_lane 1

    local.get $v
    f32x4.extract_lane 2

    f32.add
    f32.add
  )

  (func $luminance128i (param $icolor v128) (result f32)
    local.get $icolor
    f32x4.convert_i32x4_u
    call $luminance128v
  )

  (func $luminance32i (export "luminance32i")
    (param $r i32)
    (param $g i32)
    (param $b i32)
    (result f32)

    v128.const i64x2 0 0

    local.get $r
    i32x4.replace_lane 0

    local.get $g
    i32x4.replace_lane 1

    local.get $b
    i32x4.replace_lane 2

    call $luminance128i
  )

  (func $luminance32f (export "luminance32f")
    (param $r f32)
    (param $g f32)
    (param $b f32)
    (result f32)

    v128.const i64x2 0 0

    local.get $r
    f32x4.replace_lane 0

    local.get $g
    f32x4.replace_lane 1

    local.get $b
    f32x4.replace_lane 2

    call $luminance128f
  )

)
